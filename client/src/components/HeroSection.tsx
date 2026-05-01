export default function HeroSection() {
  return (
    <section className="bg-white px-4 py-20">
      <div className="page-wrap flex flex-col gap-12 md:flex-row md:items-center">
        <div className="flex-1">
          <h1 className="mb-6 text-6xl font-bold leading-tight tracking-tight">
            Kalleprint.
          </h1>
          <p className="max-w-lg text-lg text-gray-600">
            Upload your model, get an instant quote, choose your filament, and
            place your order — all in one place. Fast turnaround, quality
            prints.
          </p>
        </div>
        <div className="flex-1">
          <div className="border-2 border-black bg-gray-100 aspect-[4/3] w-full max-w-md" />
        </div>
      </div>
    </section>
  )
}
