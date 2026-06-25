import { ReactNode, forwardRef } from 'react'

interface TextAreaFieldProps {
  label: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
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
  onKeyDown,
  hint,
  testId,
  autoFocus = false,
}, ref) => {
  return (
    <div className="field">
      <div className="field-top">
        <span className="field-label">{label}</span>
        {value && (
          <button type="button" className="field-clear" onClick={() => onChange('')}>
            clear
          </button>
        )}
      </div>
      <textarea
        ref={ref}
        className="field-input"
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
