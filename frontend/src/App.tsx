import { useState, useEffect } from 'react'
import { TextAreaField } from './components/TextAreaField'
import { GenerateButton } from './components/GenerateButton'
import { StatusMessage } from './components/StatusMessage'
import { generateOrder, detectMode } from './lib/api'
import './styles/globals.css'

function App() {
  const [listMenu, setListMenu] = useState('')
  const [currentOrders, setCurrentOrders] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [status, setStatus] = useState<{ type: 'idle' | 'success' | 'error', message: string }>({
    type: 'idle',
    message: '',
  })

  const armed = Boolean(listMenu.trim() || currentOrders.trim())
  const stationLabel = armed ? detectMode(listMenu, currentOrders) : 'ready'

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
    if (isLoading) return

    // Detect mode
    const mode = detectMode(listMenu, currentOrders)

    // Validate based on mode
    if (mode === 'first-touch' && !listMenu.trim()) {
      setStatus({ type: 'error', message: 'needs a menu to print' })
      return
    }
    if ((mode === 'normal' || mode === 'nitro') && !currentOrders.trim()) {
      setStatus({ type: 'error', message: 'needs current orders to print' })
      return
    }

    setIsLoading(true)
    setStatus({ type: 'idle', message: '' })

    try {
      const data = await generateOrder({ listMenu, currentOrders, mode })

      // success — message is the raw order; the ticket renders the caption itself
      await navigator.clipboard.writeText(data.generatedMessage)
      setStatus({ type: 'success', message: data.generatedMessage })
    } catch (error) {
      setStatus({
        type: 'error',
        message: error instanceof Error ? error.message.toLowerCase() : 'something went wrong',
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

  const handleListMenuKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && e.shiftKey) {
      e.preventDefault()
      handleGenerate()
    }
  }

  return (
    <div className="ticket">
      <header className="ticket-head">
        <span className="wordmark">fast order</span>
        <span className="station" data-armed={armed}>
          <span className="station-dot" />
          {stationLabel}
        </span>
      </header>

      <TextAreaField
        label="menu"
        value={listMenu}
        onChange={setListMenu}
        placeholder="paste the menu — or leave empty for nitro"
        onKeyDown={handleListMenuKeyDown}
        testId="list-menu"
        autoFocus
      />

      <TextAreaField
        label="orders"
        value={currentOrders}
        onChange={setCurrentOrders}
        placeholder="paste what's already been ordered"
        onKeyDown={handleKeyDown}
        testId="current-orders"
      />

      <GenerateButton
        onClick={handleGenerate}
        disabled={(!listMenu.trim() && !currentOrders.trim()) || isLoading}
        loading={isLoading}
      />

      <StatusMessage type={status.type} message={status.message} testId="generated-message" />

      <div className="legend">⏎ print · ⇧⏎ first-touch · esc clear · ⌃⇧c print</div>
    </div>
  )
}

export default App
