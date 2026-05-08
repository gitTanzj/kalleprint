import type { Quote } from '#/lib/api'

interface Props {
  orderId: string
  quote: Quote
  filename: string
  customerName: string
}

export default function ConfirmationSection({ orderId, quote, filename, customerName }: Props) {
  const price = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'eur'
  }).format(quote.price_eur)

  return (
    <section className="bg-gray-700 px-4 py-20 text-white">
      <div className="page-wrap flex flex-col items-center text-center">
        <div className="mb-6 flex h-20 w-20 items-center justify-center border-2 border-white text-4xl">
          ✓
        </div>
        <h2 className="mb-3 text-4xl font-bold tracking-tight">Order Placed!</h2>
        <p className="mb-10 text-gray-300">
          Thanks, {customerName}. We'll be in touch shortly.
        </p>
        <div className="border-2 border-white p-8 text-left w-full max-w-sm">
          <p className="mb-1 text-sm uppercase tracking-widest text-gray-400">Order summary</p>
          <dl className="mt-4 space-y-3 text-sm">
            <div className="flex justify-between gap-4">
              <dt className="text-gray-400">Order ID</dt>
              <dd className="font-mono font-semibold break-all">{orderId}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt className="text-gray-400">File</dt>
              <dd className="break-all">{filename}</dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt className="text-gray-400">Price</dt>
              <dd className="font-semibold">{price}</dd>
            </div>
          </dl>
        </div>
      </div>
    </section>
  )
}
