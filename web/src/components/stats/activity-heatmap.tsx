import { useMemo } from 'react'
import { cn } from '@/lib/utils'
import type { ReviewActivity } from '@/lib/api'

interface ActivityHeatmapProps {
  data?: ReviewActivity[]
  year?: number
}

const DAY_LABELS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] as const
const VISIBLE_DAYS = new Set(['Mon', 'Wed', 'Fri'])
const MONTHS = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec',
] as const

interface DayData {
  date: string
  count: number
  isCurrentYear: boolean
}

function getIntensityClass(count: number, maxCount: number): string {
  if (count === 0) return 'bg-muted'
  const ratio = count / Math.max(maxCount, 1)
  if (ratio <= 0.25) return 'bg-primary/25'
  if (ratio <= 0.5) return 'bg-primary/50'
  if (ratio <= 0.75) return 'bg-primary/75'
  return 'bg-primary'
}

function formatDateString(date: Date): string {
  return date.toISOString().split('T')[0] ?? ''
}

export function ActivityHeatmap({ data = [], year }: ActivityHeatmapProps) {
  const currentYear = year || new Date().getFullYear()

  const { grid, maxCount, monthPositions } = useMemo(() => {
    const activityMap = new Map(data.map((d) => [d.date, d.count]))
    const startDate = new Date(currentYear, 0, 1)
    const endDate = new Date(currentYear, 11, 31)

    const firstDayOfWeek = startDate.getDay()
    const adjustedStart = new Date(startDate)
    adjustedStart.setDate(
      adjustedStart.getDate() - (firstDayOfWeek === 0 ? 6 : firstDayOfWeek - 1)
    )

    const weeks: DayData[][] = []
    const monthPos: Array<{ month: number; weekIndex: number }> = []
    const currentDate = new Date(adjustedStart)
    let max = 0
    let lastMonth = -1

    while (currentDate <= endDate || weeks.length === 0) {
      const dayOfWeek = currentDate.getDay()
      const weekIndex = dayOfWeek === 0 ? 6 : dayOfWeek - 1

      if (weekIndex === 0 || weeks.length === 0) {
        weeks.push([])
      }

      const currentWeek = weeks[weeks.length - 1]
      if (!currentWeek) continue

      const dateStr = formatDateString(currentDate)
      const count = activityMap.get(dateStr) || 0
      const isCurrentYear = currentDate.getFullYear() === currentYear

      if (count > max) max = count

      const currentMonth = currentDate.getMonth()
      if (currentMonth !== lastMonth && isCurrentYear && currentDate.getDate() <= 7) {
        monthPos.push({ month: currentMonth, weekIndex: weeks.length - 1 })
        lastMonth = currentMonth
      }

      currentWeek.push({
        date: dateStr,
        count,
        isCurrentYear,
      })

      currentDate.setDate(currentDate.getDate() + 1)

      if (currentDate > endDate && currentWeek.length === 7) {
        break
      }
    }

    return { grid: weeks, maxCount: max, monthPositions: monthPos }
  }, [data, currentYear])

  const totalReviews = data.reduce((sum, d) => sum + d.count, 0)
  const activeDays = data.filter((d) => d.count > 0).length

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>{currentYear}</span>
        <span>
          {totalReviews.toLocaleString()} reviews over {activeDays} active days
        </span>
      </div>

      <div className="overflow-x-auto">
        <div className="inline-flex flex-col gap-1">
          <div className="flex gap-1 text-xs text-muted-foreground ml-8 relative h-4">
            {monthPositions.map(({ month, weekIndex }) => (
              <span
                key={`${month}-${weekIndex}`}
                className="absolute"
                style={{ left: weekIndex * 13 }}
              >
                {MONTHS[month]}
              </span>
            ))}
          </div>
          <div className="flex gap-1 relative">
            <div className="flex flex-col gap-1 text-xs text-muted-foreground w-7 shrink-0">
              {DAY_LABELS.map((dayLabel, index) => (
                <div key={index} className="h-3 flex items-center">
                  {VISIBLE_DAYS.has(dayLabel) && <span>{dayLabel}</span>}
                </div>
              ))}
            </div>

            {grid.map((week, weekIndex) => (
              <div key={weekIndex} className="flex flex-col gap-1">
                {week.map((day) => (
                  <div
                    key={day.date}
                    className={cn(
                      'h-3 w-3 rounded-sm transition-colors',
                      day.isCurrentYear
                        ? getIntensityClass(day.count, maxCount)
                        : 'bg-transparent'
                    )}
                    title={
                      day.isCurrentYear
                        ? `${day.date}: ${day.count} ${day.count === 1 ? 'review' : 'reviews'}`
                        : undefined
                    }
                  />
                ))}
              </div>
            ))}
          </div>

          <div className="flex items-center gap-2 mt-2 text-xs text-muted-foreground">
            <span>Less</span>
            <div className="flex gap-1">
              <div className="h-3 w-3 rounded-sm bg-muted" />
              <div className="h-3 w-3 rounded-sm bg-primary/25" />
              <div className="h-3 w-3 rounded-sm bg-primary/50" />
              <div className="h-3 w-3 rounded-sm bg-primary/75" />
              <div className="h-3 w-3 rounded-sm bg-primary" />
            </div>
            <span>More</span>
          </div>
        </div>
      </div>
    </div>
  )
}
