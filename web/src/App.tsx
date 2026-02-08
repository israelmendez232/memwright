import { Toaster } from '@/components/ui/sonner'
import { MainLayout } from '@/components/layout'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'

function App() {
  return (
    <>
      <MainLayout>
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold">Welcome to Memwright</h1>
              <p className="text-muted-foreground">
                Your spaced repetition learning platform
              </p>
            </div>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Create Deck
            </Button>
          </div>

          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Cards Due Today</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">0</p>
                <p className="text-sm text-muted-foreground">
                  No cards to review
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Total Cards</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">0</p>
                <p className="text-sm text-muted-foreground">
                  Create your first deck
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Streak</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">0 days</p>
                <p className="text-sm text-muted-foreground">
                  Start studying to build your streak
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </MainLayout>
      <Toaster />
    </>
  )
}

export default App
