import { useState, useEffect } from 'react'
import { TextAreaField } from './components/TextAreaField'
import { GenerateButton } from './components/GenerateButton'
import { StatusMessage } from './components/StatusMessage'
import { generateOrder } from './lib/api'
import './styles/globals.css'

function App() {
  const [listMenu, setListMenu] = useState('')
  const [currentOrders, setCurrentOrders] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState<{ type: 'idle' | 'success' | 'error', message: string }>({
    type: 'idle',
    message: '',
  })

  // Global keyboard shortcuts
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      // ESC to clear all
      if (e.key === 'Escape') {
        setListMenu('')
        setCurrentOrders('')
        setStatus({ type: 'idle', message: '' })
      }

      // Ctrl+Shift+C to generate (when not typing)
      if (e.ctrlKey && e.shiftKey && e.key === 'C') {
        e.preventDefault()
        handleGenerate()
      }
    }

    window.addEventListener('keydown', handleGlobalKeyDown)
    return () => window.removeEventListener('keydown', handleGlobalKeyDown)
  }, [listMenu, currentOrders, isLoading])

  const handleGenerate = async () => {
    if (!currentOrders.trim()) {
      setStatus({ type: 'error', message: 'Current orders is required' })
      return
    }

    if (isLoading) return

    setIsLoading(true)
    setStatus({ type: 'idle', message: '' })

    try {
      const data = await generateOrder({ listMenu, currentOrders })

      await navigator.clipboard.writeText(data.generatedMessage)

      setStatus({ type: 'success', message: 'Order copied to clipboard! Press Ctrl+V to paste in WhatsApp' })
    } catch (error) {
      setStatus({
        type: 'error',
        message: error instanceof Error ? error.message : 'An error occurred',
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleGenerate()
    }
  }

  return (
    <div className="page-container">
      <header className="page-header">
        <h1 className="page-title">FAST ORDER</h1>
      </header>

      <main className="page-content">
        <TextAreaField
          label="LIST MENU"
          value={listMenu}
          onChange={setListMenu}
          placeholder="Paste menu here... (Ctrl+V to paste, then TAB)"
          hint="Optional - leave empty for Nitro Mode"
          testId="list-menu"
        />

        <TextAreaField
          label="CURRENT ORDERS"
          value={currentOrders}
          onChange={setCurrentOrders}
          placeholder="Paste current orders here... (Ctrl+V to paste, then ENTER to generate)"
          required
          onKeyDown={handleKeyDown}
          testId="current-orders"
        />

        <GenerateButton
          onClick={handleGenerate}
          disabled={!currentOrders.trim() || isLoading}
          loading={isLoading}
        />

        <StatusMessage type={status.type} message={status.message} testId="generated-message" />

        {status.type === 'idle' && (
          <div className="field-hint" style={{ textAlign: 'center', marginTop: 'var(--space-2)' }}>
            Shortcuts: ENTER to generate • ESC to clear • Ctrl+Shift+C to generate
          </div>
        )}
      </main>
    </div>
  )
}

export default App
