interface GenerateButtonProps {
  onClick: () => void
  disabled?: boolean
  loading?: boolean
}

export function GenerateButton({ onClick, disabled, loading }: GenerateButtonProps) {
  return (
    <div className="actions">
      <button type="button" className="print-action" onClick={onClick} disabled={disabled}>
        {loading ? 'printing…' : 'generate'}
        {!loading && <span className="kbd">⏎</span>}
      </button>
    </div>
  )
}
