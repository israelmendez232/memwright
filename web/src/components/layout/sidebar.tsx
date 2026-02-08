import { ChevronRight, FolderOpen, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

interface SidebarProps {
  className?: string
}

export function Sidebar({ className }: SidebarProps) {
  return (
    <aside
      className={cn(
        'flex w-64 flex-col border-r border-sidebar-border bg-sidebar',
        className
      )}
    >
      <div className="flex items-center justify-between p-4">
        <h2 className="text-sm font-semibold text-sidebar-foreground">Decks</h2>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          aria-label="Create new deck"
        >
          <Plus className="h-4 w-4" />
        </Button>
      </div>

      <nav className="flex-1 overflow-y-auto px-2">
        <ul className="space-y-1">
          <li>
            <a
              href="/decks/1"
              className="flex items-center gap-2 rounded-md px-3 py-2 text-sm text-sidebar-foreground transition-colors hover:bg-sidebar-accent"
            >
              <FolderOpen className="h-4 w-4 text-sidebar-primary" />
              <span className="flex-1">All Decks</span>
              <ChevronRight className="h-4 w-4 text-muted-foreground" />
            </a>
          </li>
        </ul>
      </nav>

      <div className="border-t border-sidebar-border p-4">
        <p className="text-xs text-muted-foreground">No decks yet</p>
        <p className="text-xs text-muted-foreground">
          Create your first deck to get started
        </p>
      </div>
    </aside>
  )
}
