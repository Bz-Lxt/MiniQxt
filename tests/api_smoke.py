#!/usr/bin/env python3
"""API smoke — mock/offline, expected cost ¥0."""
import json
import os
import sys
import time
import urllib.error
import urllib.request

BASE = os.environ.get("API_BASE", "http://localhost:28392")


def req(method, path, token=None, body=None, expect=None):
    data = None if body is None else json.dumps(body).encode()
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = "Bearer " + token
    r = urllib.request.Request(BASE + path, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(r, timeout=20) as resp:
            raw = resp.read().decode()
            code = resp.status
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        code = e.code
        if expect and code == expect:
            return code, json.loads(raw) if raw else {}
        raise AssertionError(f"{method} {path} -> {code} {raw}")
    payload = json.loads(raw) if raw else {}
    if expect and code != expect:
        raise AssertionError(f"{method} {path} expected {expect} got {code} {payload}")
    return code, payload


def login(user, password):
    code, body = req("POST", "/api/v1/auth/login", body={"username": user, "password": password}, expect=200)
    data = body["data"]
    return data["token"], data["user"]


def main():
    print("== health")
    code, h = req("GET", "/health", expect=200)
    assert h.get("status") == "ok", h
    print("[PASS] Health Check")

    print("== login")
    t_li, u_li = login("emp.li@hqtech", "Emp@123")
    t_chen, _ = login("emp.chen@xinghe", "Emp@123")
    t_teach, _ = login("teach.zhou@hqtech", "Teach@123")
    print("[PASS] Auth")

    print("== tenant isolation")
    # xinghe employee hitting hqtech exam 1 should 404
    try:
        req("POST", "/api/v1/exams/1/start", token=t_chen, body={})
        raise AssertionError("cross-tenant start should fail")
    except AssertionError as e:
        if "cross-tenant" in str(e):
            raise
    print("[PASS] Tenant isolation")

    print("== shuffle reproducibility")
    _, a = req("POST", "/api/v1/exams/1/start", token=t_li, body={}, expect=200)
    sid = a["data"]["session"]["id"]
    seed = a["data"]["shuffle_seed"]
    order1 = [q["question_id"] for q in a["data"]["paper"]]
    for _ in range(4):
        _, b = req("GET", f"/api/v1/exam-sessions/{sid}/paper", token=t_li, expect=200)
        order = [q["question_id"] for q in b["data"]["paper"]]
        assert order == order1, (order, order1)
        assert b["data"]["shuffle_seed"] == seed
    print("[PASS] Shuffle seed AC-07")

    print("== traces + submit")
    q0 = a["data"]["paper"][0]
    opt = q0["options"][0]["id"]
    events = []
    for i in range(12):
        events.append({
            "question_id": q0["question_id"],
            "to_option_id": opt,
            "occurred_at": time.strftime("%Y-%m-%dT%H:%M:%S+08:00"),
            "seq": i + 1,
        })
    req("POST", f"/api/v1/exam-sessions/{sid}/heartbeat", token=t_li, body={}, expect=200)
    req("POST", f"/api/v1/exam-sessions/{sid}/traces", token=t_li, body={"events": events}, expect=202)
    answers = [{"question_id": q["question_id"], "option_ids": [q["options"][0]["id"]] if q["options"] else [], "answer_text": "合规上报"} for q in a["data"]["paper"]]
    _, sub = req("POST", f"/api/v1/exam-sessions/{sid}/submit", token=t_li, body={"answers": answers}, expect=202)
    sid_sub = sub["data"]["submission_id"]
    assert sub["data"]["status"] == "queued"
    # idempotent second submit
    _, sub2 = req("POST", f"/api/v1/exam-sessions/{sid}/submit", token=t_li, body={"answers": answers}, expect=202)
    assert sub2["data"]["submission_id"] == sid_sub
    print("[PASS] Submit queue + idempotent")

    print("== poll grade")
    final = None
    for _ in range(30):
        _, g = req("GET", f"/api/v1/submissions/{sid_sub}", token=t_li, expect=200)
        st = g["data"]["submission"]["status"]
        if st in ("graded", "pending_manual"):
            final = g["data"]["submission"]
            break
        time.sleep(0.3)
    assert final, "grade timeout"
    print("[PASS] Async grade", final["status"], final["objective_score"])

    print("== integrity unverified (wang, no telemetry)")
    t_wang, _ = login("emp.wang@hqtech", "Emp@123")
    _, w = req("POST", "/api/v1/exams/1/start", token=t_wang, body={}, expect=200)
    wsid = w["data"]["session"]["id"]
    _, ws = req("POST", f"/api/v1/exam-sessions/{wsid}/submit", token=t_wang, body={"answers": []}, expect=202)
    assert ws["data"]["integrity"] == "integrity_unverified", ws
    print("[PASS] Integrity unverified without frontend telemetry")

    print("== staff metrics")
    req("GET", "/api/v1/grader/metrics", token=t_teach, expect=200)
    print("[PASS] Metrics")

    print("ALL SMOKE PASSED")
    return 0


if __name__ == "__main__":
    sys.exit(main())
