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

!insertmacro MUI_LANGUAGE "German"

; ============================================
Function .onInit

    ; --- Check 1: Is Docker installed? ---
    nsExec::ExecToStack 'cmd /C "where docker >nul 2>&1"'
    Pop $0
    Pop $1
    ${If} $0 != 0
        MessageBox MB_ICONSTOP|MB_OK "Docker wurde nicht gefunden.$\r$\n$\r$\nBitte Docker Desktop installieren und den Installer erneut starten.$\r$\nhttps://www.docker.com/products/docker-desktop"
        Abort
    ${EndIf}

    ; --- Check 2: Is Docker daemon running? ---
    nsExec::ExecToStack 'cmd /C "docker stats --no-stream >nul 2>&1"'
    Pop $0
    Pop $1
    ${If} $0 != 0
        MessageBox MB_ICONSTOP|MB_OK "Docker ist installiert, aber die Engine läuft nicht.$\r$\n$\r$\nBitte Docker Desktop starten, warten bis 'Engine running' angezeigt wird, und den Installer erneut starten."
        Abort
    ${EndIf}

FunctionEnd

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

    ${NSD_CreateLabel} 0 40u 100% 12u "Geolocation API Key:"
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
        MessageBox MB_ICONEXCLAMATION|MB_OK "Bitte eine Master Server URL eingeben."
        Abort
    ${EndIf}

    ; Validate that Geo API Key is not empty
    ${If} $GeoAPIKey == ""
        MessageBox MB_ICONEXCLAMATION|MB_OK "Bitte einen Geolocation API Key eingeben."
        Abort
    ${EndIf}

FunctionEnd

; ============================================
Section "Hauptprogramm" SecMain

    SetOutPath "$INSTDIR"

    ; Copy main executable
    File "..\..\SysTrace_Agent.exe"

    ; Create .env config file
    FileOpen $0 "$INSTDIR\.env" w
    FileWrite $0 "MASTER_SERVER_URL=$ServerURL$\r$\n"
    FileWrite $0 "GEOLOCATION_API_KEY=$GeoAPIKey$\r$\n"
    FileClose $0

    ; Install certificate
    File "..\gpshelper\GpsHelper.cer"
    ExecWait 'certutil -addstore "Root" "$INSTDIR\GpsHelper.cer"'

    ; Enable developer mode (required for sideloaded MSIX)
    WriteRegDWORD HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" "AllowAllTrustedApps" 1
    WriteRegDWORD HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" "AllowDevelopmentWithoutDevLicense" 1

    ; Register MSIX package
    File "..\gpshelper\GpsHelper.msix"
    ExecWait 'powershell -Command "Add-AppxPackage -Path \"$INSTDIR\GpsHelper.msix\""'
    Delete "$INSTDIR\GpsHelper.msix"

    ; Write uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"

    ; Start menu shortcuts
    CreateDirectory "$SMPROGRAMS\SysTrace Agent"
    CreateShortcut "$SMPROGRAMS\SysTrace Agent\SysTrace Agent.lnk" "$INSTDIR\SysTrace_Agent.exe"
    CreateShortcut "$SMPROGRAMS\SysTrace Agent\Deinstallieren.lnk" "$INSTDIR\Uninstall.exe"

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

    ; Remove GpsHelper MSIX package
    ExecWait 'powershell -Command "Get-AppxPackage -Name GpsHelper | Remove-AppxPackage"'

    ; Delete installed files
    Delete "$INSTDIR\SysTrace_Agent.exe"
    Delete "$INSTDIR\GpsHelper.cer"
    Delete "$INSTDIR\.env"
    Delete "$INSTDIR\Uninstall.exe"
    RMDir "$INSTDIR"

    ; Remove shortcuts
    Delete "$DESKTOP\SysTrace Agent.lnk"
    Delete "$SMPROGRAMS\SysTrace Agent\SysTrace Agent.lnk"
    Delete "$SMPROGRAMS\SysTrace Agent\Deinstallieren.lnk"
    RMDir "$SMPROGRAMS\SysTrace Agent"

    ; Remove registry entries
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent"
    DeleteRegKey HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock"

SectionEnd
