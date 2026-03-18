; ============================================
; SysTrace Agent - Installer
; ============================================
!include "MUI2.nsh"
!include "x64.nsh"
!include "nsDialogs.nsh"
!include "LogicLib.nsh"

Name "SysTrace Agent"
OutFile "SysTraceAgentInstaller.exe"
InstallDir "$PROGRAMFILES64\SysTrace Agent"
RequestExecutionLevel admin

!define MUI_ABORTWARNING

Var ServerURL
Var GeoAPIKey
Var Dialog
Var Label_URL
Var Label_Key
Var Input_URL
Var Input_Key

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
Page custom ConfigPageShow ConfigPageLeave
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

; ============================================
Function ConfigPageShow

    nsDialogs::Create 1018
    Pop $Dialog
    ${If} $Dialog == error
        Abort
    ${EndIf}

    ${NSD_CreateLabel} 0 0 100% 12u "Master Server URL:"
    Pop $Label_URL
    ${NSD_CreateText} 0 15u 100% 12u "http://localhost:8080"
    Pop $Input_URL

    ${NSD_CreateLabel} 0 40u 100% 12u "Geolocation API Key (ipgeolocation.io):"
    Pop $Label_Key
    ${NSD_CreateText} 0 55u 100% 12u ""
    Pop $Input_Key

    nsDialogs::Show

FunctionEnd

Function ConfigPageLeave

    ${NSD_GetText} $Input_URL $ServerURL
    ${NSD_GetText} $Input_Key $GeoAPIKey

    ; Validate that Server URL is not empty
    ${If} $ServerURL == ""
        MessageBox MB_ICONEXCLAMATION|MB_OK "Please enter a Master Server URL."
        Abort
    ${EndIf}

FunctionEnd

; ============================================
Section "Main Program" SecMain

    SetOutPath "$INSTDIR"

    ; Copy main executable
    File "..\..\SysTrace_Agent.exe"

    ; Create .env config file
    FileOpen $0 "$INSTDIR\.env" w
    FileWrite $0 "MASTER_SERVER_URL=$ServerURL$\r$\n"
    FileWrite $0 "GEOLOCATION_API_KEY=$GeoAPIKey$\r$\n"
    FileClose $0

    ; Install certificate headless
    File "..\gpshelper\GpsHelper.cer"
    nsExec::ExecToStack 'certutil -f -addstore "TrustedPeople" "$INSTDIR\GpsHelper.cer"'
    Pop $0
    Pop $1
    ${If} $0 != 0
        MessageBox MB_ICONSTOP|MB_OK "Failed to install GpsHelper certificate.$\r$\nExit code: $0$\r$\n$1"
        Abort
    ${EndIf}

    ; Enable developer mode (required for sideloaded MSIX)
    WriteRegDWORD HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" "AllowAllTrustedApps" 1
    WriteRegDWORD HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" "AllowDevelopmentWithoutDevLicense" 1

    ; Register MSIX package headless (no extra PowerShell UI)
    File "..\gpshelper\GpsHelper.msix"
    nsExec::ExecToStack 'powershell -NoLogo -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command "Add-AppxPackage -Path \"$INSTDIR\GpsHelper.msix\" -ForceApplicationShutdown"'
    Pop $0
    Pop $1
    ${If} $0 != 0
        MessageBox MB_ICONSTOP|MB_OK "Failed to install GpsHelper MSIX package."
        Abort
    ${EndIf}
    Delete "$INSTDIR\GpsHelper.msix"

    ; Write uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"

    ; Start menu shortcuts
    CreateDirectory "$SMPROGRAMS\SysTrace Agent"
    CreateShortcut "$SMPROGRAMS\SysTrace Agent\SysTrace Agent.lnk" "$INSTDIR\SysTrace_Agent.exe"
    CreateShortcut "$SMPROGRAMS\SysTrace Agent\Uninstall.lnk" "$INSTDIR\Uninstall.exe"

    ; Desktop shortcut
    CreateShortCut "$DESKTOP\SysTrace Agent.lnk" "$INSTDIR\SysTrace_Agent.exe"

    ; Add/Remove Programs entry
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "DisplayName" "SysTrace Agent"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "UninstallString" "$INSTDIR\Uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "InstallLocation" "$INSTDIR"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "DisplayVersion" "1.0.0"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "Publisher" "Elias"

SectionEnd

; ============================================
Section "Uninstall"

    ; Remove GpsHelper MSIX package headless
    nsExec::ExecToStack 'powershell -NoLogo -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command "Get-AppxPackage -Name GpsHelper | Remove-AppxPackage"'
    Pop $0
    Pop $1

    ; Delete installed files
    Delete "$INSTDIR\SysTrace_Agent.exe"
    Delete "$INSTDIR\GpsHelper.cer"
    Delete "$INSTDIR\.env"
    Delete "$INSTDIR\Uninstall.exe"
    RMDir "$INSTDIR"

    ; Remove shortcuts
    Delete "$DESKTOP\SysTrace Agent.lnk"
    Delete "$SMPROGRAMS\SysTrace Agent\SysTrace Agent.lnk"
    Delete "$SMPROGRAMS\SysTrace Agent\Uninstall.lnk"
    RMDir "$SMPROGRAMS\SysTrace Agent"

    ; Remove registry entries
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent"
    DeleteRegKey HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock"

SectionEnd
