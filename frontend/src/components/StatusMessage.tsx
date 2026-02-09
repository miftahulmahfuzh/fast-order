import { CheckCircle, AlertCircle, Loader2 } from 'lucide-react'

interface StatusMessageProps {
  type: 'idle' | 'success' | 'error' | 'loading'
  message: string
  testId?: string
}

export function StatusMessage({ type, message, testId }: StatusMessageProps) {
  if (type === 'idle') return null

  const icons = {
    success: <CheckCircle size={18} className="status-icon-success" />,
    error: <AlertCircle size={18} className="status-icon-error" />,
    loading: <Loader2 size={18} className="status-icon-loading" />,
  }

  return (
    <div className={`status-message ${type}`} data-testid={testId}>
      {icons[type]}
      {message}
    </div>
  )
}
