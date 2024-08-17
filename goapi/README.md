# go api architecture template

    ## schema oriented architecture
        - schema flow
            - request -> response (health check)
            - request -> model -> response (no via core logic)
            - request -> model -> entity -> response (via core logic)
