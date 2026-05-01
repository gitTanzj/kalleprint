import { useRef, useState, useCallback, useEffect } from 'react'
import { useFilaments } from '#/lib/queries'

const MAX_BYTES = 25 * 1024 * 1024
const ACCEPTED = ['.stl', '.3mf']

interface Props {
  onSubmit: (file: File, filamentId: string) => void
  loading: boolean
  error?: string
}

export default function UploadSection({ onSubmit, loading, error }: Props) {
  const [file, setFile] = useState<File | null>(null)
  const [fileError, setFileError] = useState<string | null>(null)
  const [dragging, setDragging] = useState(false)
  const [filamentId, setFilamentId] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const { data: filaments = [], isLoading: filamentsLoading } = useFilaments()

  function validate(f: File): string | null {
    const ext = f.name.slice(f.name.lastIndexOf('.')).toLowerCase()
    if (!ACCEPTED.includes(ext)) return 'Only .stl and .3mf files are accepted.'
    if (f.size > MAX_BYTES) return 'File must be under 25 MB.'
    return null
  }

  function pick(f: File) {
    const err = validate(f)
    if (err) {
      setFileError(err)
      setFile(null)
    } else {
      setFileError(null)
      setFile(f)
    }
  }

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDragging(false)
    const f = e.dataTransfer.files[0]
    if (f.size) pick(f)
  }, [])

  function handleSubmit() {
    if (!file || !filamentId || loading) return
    if (timerRef.current) clearTimeout(timerRef.current)
    timerRef.current = setTimeout(() => {
      onSubmit(file, filamentId)
    }, 500)
  }

  useEffect(() => () => { if (timerRef.current) clearTimeout(timerRef.current) }, [])

  const canSubmit = !!file && !!filamentId && !loading

  return (
    <section className="bg-white px-4 py-20">
      <div className="page-wrap">
        <h2 className="mb-4 text-5xl font-bold tracking-tight">Get a Quote</h2>

        <p className="mb-4">Just upload a file and get an estimate on how much the print would cost</p>

        <div
          role="button"
          tabIndex={0}
          onClick={() => inputRef.current?.click()}
          onKeyDown={(e) => e.key === 'Enter' && inputRef.current?.click()}
          onDragOver={(e) => { e.preventDefault(); setDragging(true) }}
          onDragLeave={() => setDragging(false)}
          onDrop={handleDrop}
          className={`flex min-h-56 w-full cursor-pointer flex-col items-center justify-center gap-3 border-2 border-dashed p-10 text-center transition-colors ${dragging ? 'border-black bg-gray-100' : 'border-gray-300 bg-white hover:border-black'}`}
        >
          <input
            ref={inputRef}
            type="file"
            accept=".stl,.3mf"
            className="sr-only"
            onChange={(e) => { const f = e.target.files?.[0]; if (f) pick(f) }}
          />
          {file ? (
            <>
              <span className="text-2xl">📄</span>
              <span className="font-semibold">{file.name}</span>
              <span className="text-sm text-gray-500">Click or drop to replace</span>
            </>
          ) : (
            <>
              <span className="text-2xl">⬆</span>
              <span className="font-semibold">Drop your file here</span>
              <span className="text-sm text-gray-500">or click to browse — .stl / .3mf, max 25 MB</span>
            </>
          )}
        </div>

        {fileError && <p className="mt-3 text-sm text-red-600">{fileError}</p>}

        <label className="mt-6 flex flex-col gap-1.5 text-sm font-semibold">
          Filament
          <select
            required
            value={filamentId}
            onChange={(e) => setFilamentId(e.target.value)}
            disabled={filamentsLoading}
            className="border-2 border-black bg-white px-3 py-2 font-normal text-black focus:outline-none disabled:opacity-50"
          >
            <option value="" disabled>
              {filamentsLoading ? 'Loading…' : 'Select a filament'}
            </option>
            {filaments.map((f) => (
              <option key={f.id} value={f.id}>
                {f.filament_name}
              </option>
            ))}
          </select>
        </label>

        <button
          onClick={handleSubmit}
          disabled={!canSubmit}
          className={`mt-6 border-2 border-black bg-white px-8 py-3 font-semibold text-black disabled:cursor-not-allowed cursor-pointer disabled:opacity-40${canSubmit ? ' hover-black' : ''}`}
        >
          {loading ? 'Getting quote…' : 'Get Quote'}
        </button>

        {error && <p className="mt-3 text-sm text-red-600">{error}</p>}
      </div>
    </section>
  )
}
