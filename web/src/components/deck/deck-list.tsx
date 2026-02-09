import { Loader2, Layers } from 'lucide-react'
import { useDecks } from '@/hooks'
import { DeckTree } from './deck-tree'

export function DeckList() {
  const { data: decks, isLoading, error } = useDecks()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-sm text-destructive">
        Failed to load decks. Please try again.
      </div>
    )
  }

  if (!decks || decks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
          <Layers className="h-6 w-6 text-muted-foreground" />
        </div>
        <h3 className="mt-4 font-medium">No decks yet</h3>
        <p className="mt-1 text-sm text-muted-foreground">
          Create your first deck to start learning
        </p>
      </div>
    )
  }

  return <DeckTree decks={decks} />
}
