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

test.describe('Fast Order Docker Integration', () => {
  test.use({
    baseURL: 'http://localhost:5173',
    permissions: ['clipboard-read', 'clipboard-write'],
  })

  test('should generate order from menu and current orders', async ({ page }) => {
    // Grant clipboard permissions
    await page.context().grantPermissions(['clipboard-read', 'clipboard-write'])

    await page.goto('/')

    // Fill List Menu
    await page.getByTestId('list-menu').fill(LIST_MENU)

    // Fill Current Orders
    await page.getByTestId('current-orders').fill(CURRENT_ORDERS)

    // Press Enter to generate
    await page.keyboard.press('Enter')

    // Wait for result - should show success message
    await expect(page.getByTestId('generated-message')).toBeVisible({ timeout: 30000 })

    // Verify the result contains expected patterns
    const message = await page.getByTestId('generated-message').textContent()
    // In headless mode, clipboard may fail - check for either success or clipboard error
    expect(
      message === 'copied to clipboard' || message?.includes('WhatsApp') ||
      message?.includes('Write permission denied')
    ).toBeTruthy()

    console.log('Status message:', message)
  })

  test('health check - backend is reachable', async ({ request }) => {
    const response = await request.get('http://localhost:8089/health')
    expect(response.status()).toBe(200)
  })

  test('API should generate order', async ({ request }) => {
    const response = await request.post('http://localhost:8089/api/generate-order', {
      data: {
        listMenu: LIST_MENU,
        currentOrders: CURRENT_ORDERS,
      },
    })

    expect(response.status()).toBe(200)
    const data = await response.json()
    expect(data.generatedMessage).toBeTruthy()
    expect(data.generatedMessage.length).toBeGreaterThan(50)

    console.log('Generated message preview:', data.generatedMessage.substring(0, 200))
  })
})
