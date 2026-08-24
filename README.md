# Mini 企学通

中大型企业员工培训、在线考试与岗位认证平台。完整使用说明将在 `/deploy` 阶段按交付模板补齐。

## 启动

```bash
docker compose up --build -d
```

- 员工端 http://localhost:28394  `emp.li@hqtech` / `Emp@123`
- 管理端 http://localhost:28393  `teach.zhou@hqtech` / `Teach@123`
- API http://localhost:28392/health

视频为构建期生成的示例媒体。答题流水走批量 Channel 写入；交卷先落库再入队，进程重启会 Recovery Scan。详见 `docs/Requirements.md` 与 `docs/API.md`。
