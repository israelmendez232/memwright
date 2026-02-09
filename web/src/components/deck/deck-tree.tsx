import { useState } from 'react'
import { DeckItem } from './deck-item'
import type { Deck } from '@/lib/api'

interface DeckTreeProps {
  decks: Deck[]
  level?: number
  defaultExpanded?: boolean
}

export function DeckTree({
  decks,
  level = 0,
  defaultExpanded = true,
}: DeckTreeProps) {
  const [expandedIds, setExpandedIds] = useState<Set<string>>(() => {
    if (defaultExpanded) {
      const ids = new Set<string>()
      decks.forEach((deck) => {
        if (deck.children && deck.children.length > 0) {
          ids.add(deck.id)
        }
      })
      return ids
    }
    return new Set()
  })

  const toggleExpanded = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  return (
    <div className="space-y-1">
      {decks.map((deck) => {
        const hasChildren = deck.children && deck.children.length > 0
        const isExpanded = expandedIds.has(deck.id)

        return (
          <div key={deck.id}>
            <DeckItem
              deck={deck}
              level={level}
              hasChildren={hasChildren}
              isExpanded={isExpanded}
              onToggle={() => toggleExpanded(deck.id)}
            />
            {hasChildren && isExpanded && (
              <DeckTree
                decks={deck.children!}
                level={level + 1}
                defaultExpanded={defaultExpanded}
              />
            )}
          </div>
        )
      })}
    </div>
  )
}
