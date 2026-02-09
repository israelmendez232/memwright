import { Link } from 'react-router-dom'
import { ChevronRight, ChevronDown, Layers } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { Deck } from '@/lib/api'

interface DeckItemProps {
  deck: Deck
  level?: number
  isExpanded?: boolean
  hasChildren?: boolean
  onToggle?: () => void
}

export function DeckItem({
  deck,
  level = 0,
  isExpanded = false,
  hasChildren = false,
  onToggle,
}: DeckItemProps) {
  const paddingLeft = level * 24

  return (
    <div
      className={cn(
        'group flex items-center gap-3 rounded-lg px-3 py-2.5 transition-colors hover:bg-accent/50',
        level > 0 && 'ml-3 border-l border-border'
      )}
      style={{ paddingLeft: paddingLeft + 12 }}
    >
      {hasChildren ? (
        <button
          type="button"
          onClick={(e) => {
            e.preventDefault()
            e.stopPropagation()
            onToggle?.()
          }}
          className="flex h-5 w-5 shrink-0 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground"
          aria-label={isExpanded ? 'Collapse deck' : 'Expand deck'}
        >
          {isExpanded ? (
            <ChevronDown className="h-4 w-4" />
          ) : (
            <ChevronRight className="h-4 w-4" />
          )}
        </button>
      ) : (
        <div className="w-5" />
      )}

      <Link
        to={`/decks/${deck.id}`}
        className="flex flex-1 items-center gap-3 min-w-0"
      >
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
          <Layers className="h-4 w-4" />
        </div>

        <div className="flex-1 min-w-0">
          <div className="font-medium truncate">{deck.name}</div>
          {deck.description && (
            <div className="text-sm text-muted-foreground truncate">
              {deck.description}
            </div>
          )}
        </div>

        <div className="flex shrink-0 items-center gap-4 text-sm">
          {deck.newCount > 0 && (
            <span className="flex items-center gap-1 text-blue-400">
              <span className="font-medium">{deck.newCount}</span>
              <span className="text-muted-foreground">new</span>
            </span>
          )}
          {deck.dueCount > 0 && (
            <span className="flex items-center gap-1 text-orange-400">
              <span className="font-medium">{deck.dueCount}</span>
              <span className="text-muted-foreground">due</span>
            </span>
          )}
          <span className="text-muted-foreground">
            {deck.cardCount} {deck.cardCount === 1 ? 'card' : 'cards'}
          </span>
        </div>
      </Link>
    </div>
  )
}
