import { Copy, Loader2 } from 'lucide-react'

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
      {loading ? (
        <>
          <Loader2 className="btn-icon-spin" size={20} />
          GENERATING...
        </>
      ) : (
        <>
          <Copy size={20} />
          GENERATE & COPY
        </>
      )}
    </button>
  )
}
