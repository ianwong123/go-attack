## UI evasion
When victim runs `payload.exe`, a black cmd window console appears and stays open when the binary is executed.

We can build the binary as a `GUI binary` so that the window doesn't appear with the following flag:

```
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui 
```

Or even better, we make a syscall to Windows API directly:
```
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
```

## Notes
The Windows API is a massive collection of functions that live inside DLLs sitting on disk:
1. kernel32.dll - processes, memory, files, threads
2. ntdll.dll - lowest level, talks directly to the kernel
3. user32.dll - windows, mouse, keyboard, UI
4. advapi32.dll - registry, services, security/tokens
5. ws2_32.dll - networking (WinSock)

When a `payload.exe` runs, Windows read PE (Portable Executable) headers, maps to Import Table, loads DLLs into the process's virtual memory space, resolve function address (link), and code runs