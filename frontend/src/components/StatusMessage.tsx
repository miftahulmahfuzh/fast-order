interface StatusMessageProps {
  type: 'idle' | 'success' | 'error' | 'loading'
  message: string
}

export function StatusMessage({ type, message }: StatusMessageProps) {
  if (type === 'idle') return null

  return (
    <div className={`status-message ${type}`}>
      {type === 'loading' && <div className="spinner" />}
      {message}
    </div>
  )
}
