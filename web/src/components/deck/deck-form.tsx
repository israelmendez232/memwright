import { useState, useEffect } from 'react'
import { toast } from 'sonner'
import { Loader2 } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useCreateDeck, useUpdateDeck, useDecks } from '@/hooks'
import type { Deck } from '@/lib/api'

interface FormErrors {
  name?: string
}

interface DeckFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  deck?: Deck | null
  parentId?: string | null
}

export function DeckForm({ open, onOpenChange, deck, parentId }: DeckFormProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [selectedParentId, setSelectedParentId] = useState<string | null>(null)
  const [errors, setErrors] = useState<FormErrors>({})

  const { data: decks } = useDecks()
  const createDeck = useCreateDeck()
  const updateDeck = useUpdateDeck()

  const isEditing = !!deck
  const isSubmitting = createDeck.isPending || updateDeck.isPending

  useEffect(() => {
    if (open) {
      if (deck) {
        setName(deck.name)
        setDescription(deck.description || '')
        setSelectedParentId(deck.parentId)
      } else {
        setName('')
        setDescription('')
        setSelectedParentId(parentId || null)
      }
      setErrors({})
    }
  }, [open, deck, parentId])

  const validateName = (value: string): string | undefined => {
    if (!value.trim()) {
      return 'Deck name is required'
    }
    if (value.length > 100) {
      return 'Deck name must be less than 100 characters'
    }
    return undefined
  }

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {
      name: validateName(name),
    }
    setErrors(newErrors)
    return !newErrors.name
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      if (isEditing && deck) {
        await updateDeck.mutateAsync({
          id: deck.id,
          data: {
            name: name.trim(),
            description: description.trim() || undefined,
            parentId: selectedParentId,
          },
        })
        toast.success('Deck updated successfully')
      } else {
        await createDeck.mutateAsync({
          name: name.trim(),
          description: description.trim() || undefined,
          parentId: selectedParentId,
        })
        toast.success('Deck created successfully')
      }
      onOpenChange(false)
    } catch {
      toast.error(isEditing ? 'Failed to update deck' : 'Failed to create deck')
    }
  }

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setName(value)
    if (errors.name) {
      setErrors((prev) => ({ ...prev, name: validateName(value) }))
    }
  }

  const flattenDecks = (deckList: Deck[], level = 0): Array<Deck & { level: number }> => {
    const result: Array<Deck & { level: number }> = []
    for (const d of deckList) {
      if (d.id !== deck?.id) {
        result.push({ ...d, level })
        if (d.children) {
          result.push(...flattenDecks(d.children, level + 1))
        }
      }
    }
    return result
  }

  const flatDecks = decks ? flattenDecks(decks) : []

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEditing ? 'Edit Deck' : 'Create New Deck'}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Update your deck details below.'
              : 'Add a new deck to organize your flashcards.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="name" className="text-sm font-medium">
              Name
            </label>
            <Input
              id="name"
              type="text"
              placeholder="e.g., Spanish Vocabulary"
              value={name}
              onChange={handleNameChange}
              disabled={isSubmitting}
              aria-invalid={!!errors.name}
              aria-describedby={errors.name ? 'name-error' : undefined}
            />
            {errors.name && (
              <p id="name-error" className="text-sm text-destructive">
                {errors.name}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <label htmlFor="description" className="text-sm font-medium">
              Description
              <span className="ml-1 text-muted-foreground font-normal">
                (optional)
              </span>
            </label>
            <Input
              id="description"
              type="text"
              placeholder="A brief description of this deck"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          {flatDecks.length > 0 && (
            <div className="space-y-2">
              <label htmlFor="parent" className="text-sm font-medium">
                Parent Deck
                <span className="ml-1 text-muted-foreground font-normal">
                  (optional)
                </span>
              </label>
              <select
                id="parent"
                value={selectedParentId || ''}
                onChange={(e) => setSelectedParentId(e.target.value || null)}
                disabled={isSubmitting}
                className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:border-ring disabled:cursor-not-allowed disabled:opacity-50"
              >
                <option value="">No parent (top-level deck)</option>
                {flatDecks.map((d) => (
                  <option key={d.id} value={d.id}>
                    {'  '.repeat(d.level)}
                    {d.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          <DialogFooter className="pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  {isEditing ? 'Updating...' : 'Creating...'}
                </>
              ) : isEditing ? (
                'Update Deck'
              ) : (
                'Create Deck'
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
