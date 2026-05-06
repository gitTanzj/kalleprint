import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import axios from 'axios'
import HeroSection from '#/components/HeroSection'
import UploadSection from '#/components/UploadSection'
import QuoteSection from '#/components/QuoteSection'
import ConfirmationSection from '#/components/ConfirmationSection'
import { useGetQuote } from '#/lib/queries'
import type { Quote, OrderResponse } from '#/lib/api'

export const Route = createFileRoute('/')({ component: App })

type PageState =
  | { stage: 'upload' }
  | { stage: 'quote'; quote: Quote; file: File; filamentId: string }
  | { stage: 'confirmed'; quote: Quote; file: File; filamentId: string; orderId: string; customerName: string }

function App() {
  const [state, setState] = useState<PageState>({ stage: 'upload' })
  const [networkError, setNetworkError] = useState<string | null>(null)
  const getQuote = useGetQuote()

  async function handleGetQuote(file: File, filamentId: string) {
    setNetworkError(null)
    try {
      const quote = await getQuote.mutateAsync({ file, filamentId })
      setState({ stage: 'quote', quote, file, filamentId })
    } catch (err) {
      if (axios.isAxiosError(err) && !err.response) {
        setNetworkError('Could not reach the server. Check your connection.')
      }
    }
  }

  async function handleFilamentChange(filamentId: string) {
    if (state.stage !== 'quote') return
    const file = state.file
    try {
      const quote = await getQuote.mutateAsync({ file, filamentId })
      setState((prev) => prev.stage === 'quote' ? { ...prev, quote, filamentId } : prev)
    } catch {
      // error shown via getQuote.error in QuoteSection
    }
  }

  function handleConfirmed(res: OrderResponse, customerName: string) {
    if (state.stage !== 'quote') return
    setState({ stage: 'confirmed', quote: state.quote, file: state.file, filamentId: state.filamentId, orderId: res.id, customerName })
  }

  function handleRequote() {
    getQuote.reset()
    setState({ stage: 'upload' })
  }

  return (
    <>
      {networkError && (
        <div className="border-b-2 border-black bg-red-50 px-4 py-3">
          <div className="page-wrap flex items-center justify-between gap-4">
            <p className="text-sm text-red-700">{networkError}</p>
            <button
              onClick={() => setNetworkError(null)}
              className="text-sm font-semibold text-red-700 hover:underline"
            >
              Dismiss
            </button>
          </div>
        </div>
      )}

      <HeroSection />

      {state.stage === 'upload' && (
        <UploadSection
          onSubmit={handleGetQuote}
          loading={getQuote.isPending}
          error={getQuote.error ? (getQuote.error as Error).message : undefined}
        />
      )}

      {state.stage === 'quote' && (
        <QuoteSection
          quote={state.quote}
          filename={state.file.name}
          filamentId={state.filamentId}
          isRequoting={getQuote.isPending}
          requoteError={getQuote.error?.message}
          onFilamentChange={handleFilamentChange}
          onRequote={handleRequote}
          onConfirmed={handleConfirmed}
        />
      )}

      {state.stage === 'confirmed' && (
        <ConfirmationSection
          orderId={state.orderId}
          quote={state.quote}
          filename={state.file.name}
          customerName={state.customerName}
        />
      )}
    </>
  )
}
