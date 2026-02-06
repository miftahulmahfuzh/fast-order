interface GenerateButtonProps {
  onClick: () => void
  disabled?: boolean
  loading?: boolean
}

export function GenerateButton({ onClick, disabled, loading }: GenerateButtonProps) {
  return (
    <button
      className="btn-generate"
      onClick={onClick}
      disabled={disabled}
    >
      {loading ? 'GENERATING...' : 'GENERATE & COPY'}
    </button>
  )
}
