import { Link } from 'react-router-dom'
import { BookOpen, LogOut, Settings, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context'

export function Header() {
  const { logout } = useAuth()

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-14 items-center px-4 md:px-6">
        <Link to="/" className="flex items-center gap-2">
          <BookOpen className="h-6 w-6 text-primary" />
          <span className="text-lg font-semibold">Memwright</span>
        </Link>

        <nav className="ml-8 hidden items-center gap-6 md:flex">
          <Link
            to="/decks"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Decks
          </Link>
          <Link
            to="/study"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Study
          </Link>
          <Link
            to="/stats"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Statistics
          </Link>
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <Button variant="ghost" size="icon" aria-label="Settings">
            <Settings className="h-5 w-5" />
          </Button>
          <Button variant="ghost" size="icon" aria-label="Profile">
            <User className="h-5 w-5" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            aria-label="Logout"
            onClick={() => logout()}
          >
            <LogOut className="h-5 w-5" />
          </Button>
        </div>
      </div>
    </header>
  )
}
