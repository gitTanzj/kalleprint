export default function Footer() {
  const year = new Date().getFullYear()

  return (
    <footer className="border-t-2 border-black px-4 py-8">
      <div className="page-wrap text-center text-sm text-gray-500">
        &copy; {year} Kalleprint. All rights reserved.
      </div>
    </footer>
  )
}
