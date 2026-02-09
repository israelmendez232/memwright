import { useState } from 'react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus, Flame, Calendar, Target, Clock } from 'lucide-react'
import { DeckList, DeckForm } from '@/components/deck'
import { ActivityHeatmap, StatCard } from '@/components/stats'
import { useGlobalStats, useHeatmap } from '@/hooks'

export function DashboardPage() {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const { data: stats } = useGlobalStats()
  const { data: heatmapData } = useHeatmap()

  const formatTime = (seconds: number): string => {
    if (seconds < 60) return `${seconds} secs`
    const minutes = Math.floor(seconds / 60)
    if (minutes < 60) return `${minutes} mins`
    const hours = Math.floor(minutes / 60)
    const remainingMins = minutes % 60
    return remainingMins > 0 ? `${hours}h ${remainingMins}m` : `${hours}h`
  }

  const formatPercentage = (value?: number): string => {
    if (value === undefined || value === null) return '-%'
    return `${Math.round(value)}%`
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground">
            Track your learning progress
          </p>
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Create Deck
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-blue-500/10">
                <Target className="h-5 w-5 text-blue-500" />
              </div>
              <StatCard
                label="Cards Due"
                value={stats?.cardsToReview ?? 0}
                subValue={`${stats?.newCards ?? 0} new`}
                highlight="blue"
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-green-500/10">
                <Calendar className="h-5 w-5 text-green-500" />
              </div>
              <StatCard
                label="Retention Rate"
                value={formatPercentage(stats?.retentionRate)}
                subValue="last 30 days"
                highlight="green"
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-orange-500/10">
                <Flame className="h-5 w-5 text-orange-500" />
              </div>
              <StatCard
                label="Current Streak"
                value={`${stats?.currentStreak ?? 0} days`}
                subValue={`Longest: ${stats?.longestStreak ?? 0} days`}
                highlight="orange"
              />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-purple-500/10">
                <Clock className="h-5 w-5 text-purple-500" />
              </div>
              <StatCard
                label="Review Time"
                value={formatTime(stats?.reviewTime ?? 0)}
                subValue="today"
                highlight="purple"
              />
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Review Activity</CardTitle>
          <CardDescription>
            Your daily review activity over the past year
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ActivityHeatmap data={heatmapData} />
        </CardContent>
      </Card>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Your Decks</CardTitle>
                <CardDescription>
                  {stats?.totalDecks ?? 0} decks with {stats?.totalCards ?? 0} cards
                </CardDescription>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsCreateDialogOpen(true)}
              >
                <Plus className="mr-2 h-4 w-4" />
                New Deck
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <DeckList />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Learning Stats</CardTitle>
            <CardDescription>Your overall progress</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Daily Average</span>
                <span className="font-medium text-blue-400">
                  {stats?.dailyAverage ?? 0} cards
                </span>
              </div>
              <div className="h-2 rounded-full bg-muted">
                <div
                  className="h-2 rounded-full bg-blue-500"
                  style={{
                    width: `${Math.min((stats?.dailyAverage ?? 0) / 100 * 100, 100)}%`,
                  }}
                />
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Days Learned</span>
                <span className="font-medium text-green-400">
                  {stats?.totalDays
                    ? `${Math.round((stats.daysLearned / stats.totalDays) * 100)}%`
                    : '0%'}
                </span>
              </div>
              <div className="h-2 rounded-full bg-muted">
                <div
                  className="h-2 rounded-full bg-green-500"
                  style={{
                    width: `${
                      stats?.totalDays
                        ? (stats.daysLearned / stats.totalDays) * 100
                        : 0
                    }%`,
                  }}
                />
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Cards Learned</span>
                <span className="font-medium text-purple-400">
                  {stats?.cardsLearned ?? 0} / {stats?.totalCards ?? 0}
                </span>
              </div>
              <div className="h-2 rounded-full bg-muted">
                <div
                  className="h-2 rounded-full bg-purple-500"
                  style={{
                    width: `${
                      stats?.totalCards
                        ? (stats.cardsLearned / stats.totalCards) * 100
                        : 0
                    }%`,
                  }}
                />
              </div>
            </div>

            <div className="pt-4 border-t">
              <div className="grid grid-cols-2 gap-4 text-center">
                <div>
                  <p className="text-2xl font-bold">{stats?.totalReviews ?? 0}</p>
                  <p className="text-xs text-muted-foreground">Total Reviews</p>
                </div>
                <div>
                  <p className="text-2xl font-bold">{stats?.masteredCards ?? 0}</p>
                  <p className="text-xs text-muted-foreground">Mastered</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <DeckForm open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen} />
    </div>
  )
}
