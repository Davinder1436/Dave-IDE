flowchart LR
    subgraph UI["User Interface"]
        U1[User Login]
        U2[Room Selection]
        U3[IDE Interaction]
    end
    
    subgraph Backend["Backend API"]
        B1[User Auth Service]
        B2[Room Manager]
        B3[Container Manager]
        B4[Git Service]
        B5[S3 Storage Manager]
    end
    
    subgraph ECS["ECS Service"]
        E1[Create Container]
        E2[Manage Container Lifecycle]
    end
    
    subgraph Docker["Docker Container"]
        D1[IDE Environment]
        D2[Fixed Size Volume]
    end
    
    UI -->|Login Request| B1
    B1 -->|Auth Response| UI
    UI -->|Room Request| B2
    B2 -->|Create Room| B3
    B3 -->|Create Task| ECS
    ECS -->|Container Created| Docker
    Docker -->|Room Access| UI
    UI -->|Code Interaction| Docker
    Docker -->|File Updates| B4
    B4 -->|Git Commit| S3
    B5 -->|Store Backup| S3
