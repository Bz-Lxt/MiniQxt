import { test, expect } from '@playwright/test'

const USER = 'http://localhost:28394'
const ADMIN = 'http://localhost:28393'

test('employee login learn exam result', async ({ page }) => {
  await page.goto(USER + '/login')
  await page.fill('input', 'emp.li@hqtech')
  await page.locator('input[type="password"]').fill('Emp@123')
  await page.getByRole('button', { name: '进入学习' }).click()
  await expect(page.getByText('待学课程')).toBeVisible()
  await page.getByText('进入工作台').first().click()
  await expect(page.getByText('有效进度')).toBeVisible()
  await page.goto(USER + '/')
  await page.getByText('进入考场').first().click()
  await expect(page.getByText('IMMERSIVE EXAM')).toBeVisible()
})

test('admin paper board', async ({ page }) => {
  await page.goto(ADMIN + '/login')
  await page.locator('input').first().fill('teach.zhou@hqtech')
  await page.locator('input[type="password"]').fill('Teach@123')
  await page.getByRole('button', { name: '进入台账' }).click()
  await expect(page.getByText('夜间台账')).toBeVisible()
  await page.getByText('出题看板').click()
  await expect(page.getByText('拖拽智能出题看板')).toBeVisible()
})
