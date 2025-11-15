# Velock-Web Architecture Documentation

## System Overview

Velock-web is a web application for desktop app authentication and friend invitation management, built with React and Firebase.

## UML Diagrams

### System Architecture Diagram

```mermaid
graph TB
    subgraph "Client Layer"
        Browser[Web Browser<br/>localhost:4444]
        DesktopApp[Desktop App<br/>Custom Protocol]
    end

    subgraph "Frontend - React App"
        Router[React Router]
        Store[Zustand Store<br/>Persistent State]
        AuthListener[Auth Listener]

        subgraph "Components"
            DesktopSignIn[Desktop Sign In]
            InviteScreen[Invite Screen]
            GoogleBtn[Google Sign In Button]
        end
    end

    subgraph "Firebase Backend"
        subgraph "Firebase Auth"
            GoogleAuth[Google OAuth]
            CustomAuth[Custom Tokens]
        end

        subgraph "Cloud Functions"
            CreateToken[createAuthToken]
            PostSignup[postSignupSteps]
            CreateInvite[createInvite]
            AcceptInvite[acceptInvite]
        end

        subgraph "Firestore"
            UserDetails[(user_details)]
            Invites[(invites)]
            Registry[(registry)]
        end

        subgraph "Realtime DB"
            AuthCodes[(ot-auth-codes)]
        end
    end

    subgraph "External Services"
        PostHog[PostHog Analytics]
        EmailService[Email Service<br/>TODO: Loops/Brevo]
    end

    Browser --> Router
    Router --> DesktopSignIn
    Router --> InviteScreen

    DesktopSignIn --> GoogleBtn
    InviteScreen --> GoogleBtn

    GoogleBtn --> GoogleAuth
    GoogleAuth --> AuthListener
    AuthListener --> Store

    DesktopSignIn --> CreateToken
    InviteScreen --> AcceptInvite

    CreateToken --> CustomAuth
    CreateToken --> AuthCodes

    PostSignup --> UserDetails
    CreateInvite --> Invites
    CreateInvite --> EmailService
    AcceptInvite --> Invites
    AcceptInvite --> Registry

    DesktopApp -.->|Opens| Browser
    Browser -.->|Redirects| DesktopApp

    Browser --> PostHog
```

### Data Model Diagram

```mermaid
classDiagram
    class User {
        +string uid
        +string email
        +string name
        +string photoURL
    }

    class InviteData {
        +string id
        +string inviterUid
        +string inviteeEmail
        +InviteMeta meta
        +string createdAt
        +string status
    }

    class InviteMeta {
        +UserInfo inviter
        +UserInfo invitee
    }

    class UserInfo {
        +string uid
        +string name
        +string email
        +string photoURL
    }

    class FriendDetails {
        +string uid
        +string name
        +string email
        +string photoURL
    }

    class Registry {
        +string userId
        +Map~string,FriendDetails~ friends
    }

    class AppStore {
        +boolean isLoggedIn
        +User user
        +setUserLoggedIn()
        +setUser()
    }

    InviteData --> InviteMeta
    InviteMeta --> UserInfo
    Registry --> FriendDetails
    AppStore --> User
```

### Component Hierarchy Diagram

```mermaid
graph TD
    Main[main.tsx<br/>Entry Point] --> RouterProvider[RouterProvider]
    RouterProvider --> App[App.tsx]

    App --> AuthListener[AuthListener<br/>Global Auth Sync]
    App --> Outlet[Outlet<br/>Route Content]

    Outlet --> DesktopSignIn[DesktopSignIn<br/>/desktop-login]
    Outlet --> InviteScreen[InviteScreen<br/>/invite/:id]
    Outlet --> NotFound[NotFound<br/>/*]

    DesktopSignIn --> OnboardingCard[OnboardingCard]
    OnboardingCard --> GoogleSignInButton1[GoogleSignInButton]
    OnboardingCard --> AuthMessageCard[AuthMessageCard]

    InviteScreen --> GoogleSignInButton2[GoogleSignInButton]
    InviteScreen --> PageLoader[PageLoader]

    AuthListener -.->|Updates| Store[Zustand Store]
    GoogleSignInButton1 -.->|Triggers| FirebaseAuth[Firebase Auth]
    GoogleSignInButton2 -.->|Triggers| FirebaseAuth

    style Store fill:#f9f,stroke:#333,stroke-width:2px
    style FirebaseAuth fill:#ff9,stroke:#333,stroke-width:2px
```

### Authentication Flow Diagram

```mermaid
sequenceDiagram
    participant Desktop as Desktop App
    participant Browser as Web Browser
    participant Firebase as Firebase Auth
    participant Function as Cloud Function
    participant DB as Realtime DB
    participant Store as Zustand Store

    Desktop->>Browser: Open with one-time code
    Browser->>Browser: Display Google Sign In
    Browser->>Firebase: signInWithPopup(Google)
    Firebase-->>Browser: ID Token

    Firebase->>Store: onAuthStateChanged
    Store->>Store: Update user state

    Browser->>Function: createAuthToken(code, idToken)
    Function->>Firebase: Verify ID Token
    Function->>Function: Create custom token
    Function->>DB: Store token at ot-auth-codes/{code}
    Function-->>Browser: Success

    Browser->>Desktop: Redirect via pairing-client://
    Desktop->>DB: Poll for token
    DB-->>Desktop: Custom token
    Desktop->>Firebase: signInWithCustomToken
```

### Invite Flow Diagram

```mermaid
sequenceDiagram
    participant Inviter as Inviter (Desktop)
    participant Function1 as createInvite
    participant Firestore as Firestore
    participant Email as Email Service
    participant Browser as Invitee (Browser)
    participant Function2 as acceptInvite

    Inviter->>Function1: createInvite(email)
    Function1->>Function1: Classify email type
    Function1->>Function1: Check existing invite
    Function1->>Firestore: Create invite document
    Function1->>Email: Send invite email
    Function1-->>Inviter: Success + inviteId

    Email->>Browser: Invitee clicks link
    Browser->>Browser: Sign in with Google
    Browser->>Firestore: Fetch invite details
    Firestore-->>Browser: Invite data
    Browser->>Browser: Display inviter info
    Browser->>Function2: acceptInvite(inviteId)
    Function2->>Function2: Verify invitee email
    Function2->>Firestore: Create bidirectional friendship
    Function2->>Firestore: Update invite status
    Function2-->>Browser: Success
    Browser->>Browser: Show success message
```

### Security Rules Diagram

```mermaid
graph LR
    subgraph "Firestore Collections"
        Users[(users)]
        Registry[(registry)]
        Invites[(invites)]
    end

    subgraph "Access Rules"
        UsersRule["Read/Update:<br/>Owner only<br/>Create: Functions only"]
        RegistryRule["Read/Update:<br/>Owner only<br/>Create: Functions only"]
        InvitesRule["Read:<br/>Inviter OR Invitee email<br/>Update: Inviter only<br/>Create: Functions only<br/>List: Any authenticated"]
    end

    subgraph "Realtime DB"
        AuthCodes[(ot-auth-codes)]
        OpenRule["Read/Write:<br/>Anyone<br/>⚠️ INSECURE"]
    end

    Users --> UsersRule
    Registry --> RegistryRule
    Invites --> InvitesRule
    AuthCodes --> OpenRule

    style OpenRule fill:#f66,stroke:#333,stroke-width:2px
```

## Architecture Layers

### 1. Client Layer
- **Web Browser**: React SPA running on Vite dev server (port 4444)
- **Desktop App**: Native app using custom URL protocol (`pairing-client://`)

### 2. Frontend Layer
- **Framework**: React 19 + Vite
- **Routing**: React Router 7
- **State Management**: Zustand with localStorage persistence
- **Styling**: TailwindCSS 4
- **Analytics**: PostHog

### 3. Backend Layer
- **Firebase Auth**: Google OAuth + Custom tokens
- **Cloud Functions**: 4 callable functions + 1 auth trigger
- **Firestore**: 3 collections (user_details, invites, registry)
- **Realtime Database**: One-time auth codes

### 4. External Services
- **PostHog**: User analytics and event tracking
- **Email Service**: TODO - Loops or Brevo integration

## Key Design Patterns

### 1. Global State Management
- **Pattern**: Zustand store with localStorage persistence
- **Purpose**: Maintain auth state across page reloads
- **Implementation**: `app/src/modules/store.ts`

### 2. Auth State Synchronization
- **Pattern**: Observer pattern via Firebase onAuthStateChanged
- **Purpose**: Sync Firebase auth with local Zustand state
- **Implementation**: `AuthListener` component

### 3. Security by Default
- **Pattern**: Deny-all Firestore rules with explicit allows
- **Purpose**: Prevent unauthorized data access
- **Implementation**: `cloud/firestore.rules`

### 4. Email Classification
- **Pattern**: Strategy pattern for different email types
- **Purpose**: Handle invites differently based on email domain
- **Implementation**: `cloud/functions/src/utils/email.ts`

### 5. Atomic Friendship
- **Pattern**: Firestore transaction for bidirectional updates
- **Purpose**: Ensure both users become friends simultaneously
- **Implementation**: `makeFriends()` in `cloud/functions/src/utils/friend.ts`

## Data Flow

### Authentication Flow
1. User clicks Google Sign In
2. Firebase Auth popup opens
3. User authorizes
4. Firebase returns ID token
5. `onAuthStateChanged` triggers
6. `AuthListener` updates Zustand store
7. Components re-render with auth state

### Desktop Auth Flow
1. Desktop app opens browser with one-time code
2. User authenticates via Google
3. Web calls `createAuthToken(code, idToken)`
4. Function creates custom token
5. Token stored in Realtime DB
6. Web redirects to desktop app
7. Desktop app polls Realtime DB
8. Desktop app uses custom token to authenticate

### Invite Flow
1. User A creates invite via desktop app
2. `createInvite(email)` called
3. Function classifies email type
4. Invite document created in Firestore
5. Email sent to invitee
6. User B clicks invite link
7. User B authenticates
8. Invite details fetched (security rules check email)
9. User B clicks accept
10. `acceptInvite(inviteId)` called
11. Function verifies email match
12. Bidirectional friendship created
13. Invite status updated to 'accepted'

## Security Considerations

### ✅ Secure
- Firestore rules enforce owner-only access
- Invite reads require email match or inviter UID
- All writes go through authenticated Cloud Functions
- Firebase Auth handles authentication

### ⚠️ Needs Attention
- **Realtime Database**: Open read/write (CRITICAL - fix before production)
- **Email validation**: Should validate email format in Cloud Functions
- **Rate limiting**: No rate limiting on invite creation
- **Input sanitization**: Limited input validation

## Performance Optimizations

1. **Zustand shallow comparison**: Prevents unnecessary re-renders
2. **Firebase emulators**: Fast local development
3. **Vite HMR**: Hot module replacement for fast dev feedback
4. **PostHog batching**: Attributes synced every 20s

## Deployment Architecture

### Development
- Frontend: `npm run dev:app` (Vite on port 4444)
- Backend: `npm run dev:cloud` (Firebase emulators)
- Combined: `npm run dev` (Concurrently)

### Production
- Frontend: Deployed to hosting (TBD)
- Backend: Firebase Cloud Functions (auto-scaling)
- Database: Firestore + Realtime DB (managed by Firebase)

## Scalability Considerations

1. **Cloud Functions**: Max 10 instances (configurable)
2. **Firestore**: Scales automatically
3. **Realtime DB**: Single instance, may need sharding at scale
4. **Frontend**: Static hosting, scales with CDN

## Future Enhancements

### High Priority
1. Secure Realtime Database rules
2. Implement email service integration
3. Add invite link generation with proper URLs
4. Auto-accept pending invites on signup

### Medium Priority
5. Implement pairing discovery service
6. Add rate limiting to Cloud Functions
7. Improve error handling and user feedback
8. Add comprehensive logging

### Low Priority
9. Expand analytics tracking
10. Add unit and integration tests
11. Implement invite expiration
12. Add friend removal functionality
