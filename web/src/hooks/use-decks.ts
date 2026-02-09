import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  deckApi,
  statsApi,
  type Deck,
  type CreateDeckRequest,
  type UpdateDeckRequest,
} from '@/lib/api'

const DECK_KEYS = {
  all: ['decks'] as const,
  lists: () => [...DECK_KEYS.all, 'list'] as const,
  list: (filters?: Record<string, unknown>) =>
    [...DECK_KEYS.lists(), filters] as const,
  details: () => [...DECK_KEYS.all, 'detail'] as const,
  detail: (id: string) => [...DECK_KEYS.details(), id] as const,
  stats: (id: string) => [...DECK_KEYS.all, 'stats', id] as const,
}

const STATS_KEYS = {
  all: ['stats'] as const,
  global: () => [...STATS_KEYS.all, 'global'] as const,
  heatmap: (year?: number) => [...STATS_KEYS.all, 'heatmap', year] as const,
}

function buildDeckTree(decks: Deck[]): Deck[] {
  const deckMap = new Map<string, Deck>()
  const rootDecks: Deck[] = []

  decks.forEach((deck) => {
    deckMap.set(deck.id, { ...deck, children: [] })
  })

  deckMap.forEach((deck) => {
    if (deck.parentId && deckMap.has(deck.parentId)) {
      const parent = deckMap.get(deck.parentId)!
      parent.children = parent.children || []
      parent.children.push(deck)
    } else {
      rootDecks.push(deck)
    }
  })

  return rootDecks
}

export function useDecks() {
  return useQuery({
    queryKey: DECK_KEYS.lists(),
    queryFn: async () => {
      const response = await deckApi.getAll()
      return buildDeckTree(response.data)
    },
  })
}

export function useDeck(id: string) {
  return useQuery({
    queryKey: DECK_KEYS.detail(id),
    queryFn: async () => {
      const response = await deckApi.getById(id)
      return response.data
    },
    enabled: !!id,
  })
}

export function useDeckStats(id: string) {
  return useQuery({
    queryKey: DECK_KEYS.stats(id),
    queryFn: async () => {
      const response = await deckApi.getStats(id)
      return response.data
    },
    enabled: !!id,
  })
}

export function useCreateDeck() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: CreateDeckRequest) => {
      const response = await deckApi.create(data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DECK_KEYS.lists() })
    },
  })
}

export function useUpdateDeck() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateDeckRequest }) => {
      const response = await deckApi.update(id, data)
      return response.data
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: DECK_KEYS.lists() })
      queryClient.invalidateQueries({ queryKey: DECK_KEYS.detail(variables.id) })
    },
  })
}

export function useDeleteDeck() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      await deckApi.delete(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: DECK_KEYS.lists() })
    },
  })
}

export function useGlobalStats() {
  return useQuery({
    queryKey: STATS_KEYS.global(),
    queryFn: async () => {
      const response = await statsApi.getGlobal()
      return response.data
    },
  })
}

export function useHeatmap(year?: number) {
  return useQuery({
    queryKey: STATS_KEYS.heatmap(year),
    queryFn: async () => {
      const response = await statsApi.getHeatmap(year)
      return response.data
    },
  })
}
