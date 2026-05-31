import { useRef, useState, useCallback, useEffect } from 'react'
import { useFilaments } from '#/lib/queries'
import type { Filament } from '#/lib/api'

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
  const [dropdownOpen, setDropdownOpen] = useState(false)
  const [hoveredFilament, setHoveredFilament] = useState<Filament | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)

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

  useEffect(() => {
    if (!dropdownOpen) return
    function handleClickOutside(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [dropdownOpen])

  const selectedFilament = filaments.find((f) => f.id === filamentId);

  function selectFilament(id: string) {
    setFilamentId(id)
    setDropdownOpen(false)
  }

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

        <div className="mt-6 flex flex-col gap-1.5 text-sm font-semibold">
          Filament
          <div ref={dropdownRef} className="relative">
            <button
              type="button"
              onClick={() => !filamentsLoading && setDropdownOpen((o) => !o)}
              disabled={filamentsLoading}
              className="flex w-full items-center gap-2.5 border-2 border-black bg-white px-3 py-2 text-left font-normal text-black focus:outline-none disabled:opacity-50"
            >
              {selectedFilament ? (
                <>
                  <span
                    className="inline-block h-4 w-4 shrink-0 rounded-full border border-black/20"
                    style={{ backgroundColor: `#${selectedFilament.color_hex}` }}
                  />
                  <span className="flex-1">{selectedFilament.filament_name}</span>
                </>
              ) : (
                <span className="flex-1 text-gray-400">
                  {filamentsLoading ? 'Loading…' : 'Select a filament'}
                </span>
              )}
              <svg
                className={`h-4 w-4 shrink-0 transition-transform ${dropdownOpen ? 'rotate-180' : ''}`}
                viewBox="0 0 16 16"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
              >
                <path d="M4 6l4 4 4-4" />
              </svg>
            </button>

            {dropdownOpen && hoveredFilament && (
              <div className="pointer-events-none absolute top-0 right-[calc(100%+12px)] flex flex-col items-center gap-1.5 z-20">
                <div
                  className="h-16 w-16 border-2 border-black"
                  style={{ backgroundColor: `#${hoveredFilament.color_hex}` }}
                />
                <span className="font-mono text-xs">{hoveredFilament.color_hex}</span>
              </div>
            )}

            {dropdownOpen && (
              <ul className="absolute z-10 mt-[-2px] max-h-60 w-full overflow-auto border-2 border-black bg-white">
                {filaments.map((f) => (
                  <li key={f.id}>
                    <button
                      type="button"
                      onClick={() => selectFilament(f.id)}
                      className={`flex w-full items-center gap-2.5 px-3 py-2 text-left font-normal hover:bg-gray-100 ${f.id === filamentId ? 'bg-gray-100 font-semibold' : ''}`}
                    >
                      <span
                        className="inline-block h-4 w-4 shrink-0 rounded-full border border-black/20 cursor-default"
                        style={{ backgroundColor: `#${f.color_hex}` }}
                        onMouseEnter={() => setHoveredFilament(f)}
                        onMouseLeave={() => setHoveredFilament(null)}
                      />
                      {f.filament_name}
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

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
