# Expo Router App

This is an Expo app using Expo Router for navigation with a comprehensive dashboard featuring tabs.

## Structure

- `app/` - Contains all the routes and layouts
  - `_layout.tsx` - Root layout that defines the navigation structure
  - `index.tsx` - Home page (default route)
  - `settings.tsx` - Settings page
  - `dashboard/` - Dashboard with tabs navigation
    - `_layout.tsx` - Tab layout for dashboard
    - `overview.tsx` - Dashboard overview tab
    - `expenses.tsx` - Expenses tracking tab
    - `revenue.tsx` - Revenue tracking tab
    - `goals.tsx` - Financial goals tab

## Navigation

The app uses Expo Router for file-based routing. Each file in the `app/` directory becomes a route:

- `/` - Home page (index.tsx)
- `/dashboard` - Dashboard with tabs (overview, expenses, revenue, goals)
- `/settings` - Settings page (settings.tsx)

### Dashboard Tabs

The dashboard features four main tabs:

1. **Overview** - Financial summary and recent activity
2. **Expenses** - Expense tracking and categories
3. **Revenue** - Income tracking and sources
4. **Goals** - Financial goal management

## Adding New Routes

To add a new route, simply create a new file in the `app/` directory:

```tsx
// app/profile.tsx
export default function Profile() {
  return <Text>Profile Page</Text>;
}
```

Then add it to the Stack in `_layout.tsx`:

```tsx
<Stack.Screen name="profile" options={{ title: 'Profile' }} />
```

## Adding New Dashboard Tabs

To add a new tab to the dashboard:

1. Create a new file in `app/dashboard/` (e.g., `app/dashboard/analytics.tsx`)
2. Add it to the Tabs in `app/dashboard/_layout.tsx`:

```tsx
<Tabs.Screen
  name="analytics"
  options={{
    title: 'Analytics',
    tabBarIcon: ({ color, size }) => (
      <Ionicons name="bar-chart" size={size} color={color} />
    ),
  }}
/>
```

## Running the App

```bash
npm start
```

## Building for Android

### Development Build
```bash
npm run build:dev:android
```

### Preview Build
```bash
npm run build:preview:android
```

### Production Build
```bash
npm run build:prod:android
```

For detailed build instructions, see [BUILD_GUIDE.md](./BUILD_GUIDE.md).

## Features

- File-based routing with Expo Router
- Tab navigation in dashboard
- TypeScript support
- NativeWind for styling
- Tailwind CSS
- Ionicons for icons
- Financial tracking UI components
- Import aliases for cleaner imports
- EAS Build configuration for Android

## Import Aliases

The project uses import aliases to make imports cleaner and more maintainable:

- `@/*` - Root directory imports
- `~/components/*` - Component imports
- `~/assets/*` - Asset imports
- `~/app/*` - App route imports

### Examples:

```tsx
// Before
import { ScreenContent } from '../../components/ScreenContent';
import '../global.css';

// After
import { ScreenContent } from '@/components/ScreenContent';
import '@/global.css';
``` 