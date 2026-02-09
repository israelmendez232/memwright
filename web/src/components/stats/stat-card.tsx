import { cn } from '@/lib/utils'

interface StatCardProps {
  label: string
  value: string | number
  subValue?: string
  className?: string
  highlight?: 'blue' | 'green' | 'orange' | 'purple'
}

export function StatCard({
  label,
  value,
  subValue,
  className,
  highlight,
}: StatCardProps) {
  const highlightClass = {
    blue: 'text-blue-400',
    green: 'text-green-400',
    orange: 'text-orange-400',
    purple: 'text-purple-400',
  }

  return (
    <div className={cn('space-y-1', className)}>
      <p className="text-sm text-muted-foreground">{label}</p>
      <p
        className={cn(
          'text-2xl font-semibold',
          highlight && highlightClass[highlight]
        )}
      >
        {value}
      </p>
      {subValue && (
        <p className="text-xs text-muted-foreground">{subValue}</p>
      )}
    </div>
  )
}
