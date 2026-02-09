import { test, expect } from '@playwright/test'

const LIST_MENU = `Menu Redi now 🥰🥰
Cah buncis
Cah labuu
Cah terong balado
Cah jagung muda
Cah toge
Kari kentanggg
Usus cabe ijo
Kikil setan
Ati ampela tahu cabe ijo
Ceker setan
Oncom leuncha
Tempe lada hitamsss
Tempe sambal pecaks
Tahu sambal pecaks
Jamur crispyyy
Ayam suwir daun jeruk
Sup bakso sapi
Baby cumi cabe ijo
Fillet ayam crispy cabe garam
Udang crispy cabe ijo
Ikan mujair bakarss
Ayam betutu baliii
Dendeng sapi cabe ijo
Telur bulet bumbu kuning
Orek basah
Kentang balado
Telur ceplok gulai
Perkedel
Rolade ayam
Sosis oseng telur
Otak otak cb garem
Mie goreng
Bihun goreng
Tempe cabe garam
Pepes tahu


Gorengan ; tempe🥰

Es mambo
Kacang ijo , sirup 🥰

Susu kacang
Original, coklat 🥰

Jajanan kue pasar🥰
Donut, lemper, dadar kelapa , lapis Surabaya, onde, gemblong, pastel , risol sayur , risol mayo`

const CURRENT_ORDERS = `06/02/26

1. farid : nasi 1/2 + cah buncis + ati ampela cabe ijo + jamur crispy + mie goreng
2. Rian TII: Nasi 1/2 + Cah Jangung Muda + Tempe Lada Hitam + Perkedel + Fillet Ayam Cabe Garam + Sambal + kuah sedikit
3.`

test.describe('Fast Order Integration', () => {
  test('should generate order from menu and current orders', async ({ page }) => {
    await page.goto('/')

    // Fill List Menu
    await page.getByTestId('list-menu').fill(LIST_MENU)

    // Fill Current Orders
    await page.getByTestId('current-orders').fill(CURRENT_ORDERS)

    // Press Enter to generate
    await page.keyboard.press('Enter')

    // Wait for result - should show generated message
    await expect(page.getByTestId('generated-message')).toBeVisible({ timeout: 30000 })

    // Verify the result contains expected patterns
    const message = await page.getByTestId('generated-message').textContent()
    expect(message).toBeTruthy()
    expect(message!.length).toBeGreaterThan(50)

    console.log('Generated message:', message)
  })

  test('health check - backend is reachable', async ({ request }) => {
    const response = await request.get('http://localhost:8089/health')
    expect(response.status()).toBe(200)
  })
})
