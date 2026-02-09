import { ReactNode } from 'react'

interface TextAreaFieldProps {
  label: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  onKeyDown?: (e: React.KeyboardEvent) => void
  hint?: ReactNode
  testId?: string
}

export function TextAreaField({
  label,
  value,
  onChange,
  placeholder,
  required = false,
  onKeyDown,
  hint,
  testId,
}: TextAreaFieldProps) {
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
          >
            CLEAR
          </button>
        )}
      </div>
      <textarea
        className="brutalist-textarea"
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={onKeyDown}
        data-testid={testId}
      />
      {hint && <div className="field-hint">{hint}</div>}
    </div>
  )
}
