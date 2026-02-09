import { ReactNode, forwardRef } from 'react'
import { X } from 'lucide-react'

interface TextAreaFieldProps {
  label: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  onKeyDown?: (e: React.KeyboardEvent) => void
  hint?: ReactNode
  testId?: string
  autoFocus?: boolean
}

export const TextAreaField = forwardRef<HTMLTextAreaElement, TextAreaFieldProps>(({
  label,
  value,
  onChange,
  placeholder,
  required = false,
  onKeyDown,
  hint,
  testId,
  autoFocus = false,
}, ref) => {
  return (
    <div className="field-container">
      <div className="field-header">
        <span className="field-label">
          {label}
          {required && <span className="field-required"> *</span>}
        </span>
        {value && (
          <button
            className="btn-clear"
            onClick={() => onChange('')}
            aria-label="Clear"
          >
            <X size={16} />
          </button>
        )}
      </div>
      <textarea
        ref={ref}
        className="brutalist-textarea"
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={onKeyDown}
        autoFocus={autoFocus}
        data-testid={testId}
      />
      {hint && <div className="field-hint">{hint}</div>}
    </div>
  )
})

TextAreaField.displayName = 'TextAreaField'
