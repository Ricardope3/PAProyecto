# Threads
- main thread executes the run method
- the run method executes each person's goroutine
- the run method executes the select statement within an anonimous goroutine
- each person communicates through the `onMove` and `onExit` channels to the select statement
- the run method keeps executing an update method to re-render the canvas
![](assets/threads.png)
# METHODS
