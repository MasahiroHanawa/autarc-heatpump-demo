import { Link } from 'react-router-dom'

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold text-gray-900">
            Heat Pump Planning Assistant
          </h1>
          <nav className="flex items-center gap-4">
            <Link to="/" className="text-sm text-gray-600 hover:text-gray-900">
              New Analysis
            </Link>
            <Link to="/jobs" className="text-sm text-gray-600 hover:text-gray-900">
              History
            </Link>
          </nav>
        </div>
      </header>
      <main className="max-w-4xl mx-auto px-4 py-8">{children}</main>
    </div>
  )
}
