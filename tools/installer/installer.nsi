; ============================================
; SysTrace Agent - Installer
; ============================================
!include "MUI2.nsh"
!include "x64.nsh"

Name "SysTrace Agent"
OutFile "SysTraceAgentInstaller.exe"
InstallDir "$PROGRAMFILES64\SysTrace Agent"
RequestExecutionLevel admin

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "German"

; ============================================
Section "Hauptprogramm" SecMain

    SetOutPath "$INSTDIR"

    File "..\..\SysTrace_Agent.exe"

    File "..\gpshelper\GpsHelper.cer"
    ExecWait 'certutil -addstore "Root" "$INSTDIR\GpsHelper.cer"'

    WriteRegDWORD HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" "AllowAllTrustedApps" 1
    WriteRegDWORD HKLM "SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock" "AllowDevelopmentWithoutDevLicense" 1

    File "..\gpshelper\GpsHelper.msix"
    ExecWait 'powershell -Command "Add-AppxPackage -Path \"$INSTDIR\GpsHelper.msix\""'
    Delete "$INSTDIR\GpsHelper.msix"

    WriteUninstaller "$INSTDIR\Uninstall.exe"

    CreateDirectory "$SMPROGRAMS\SysTrace Agent"
    CreateShortcut "$SMPROGRAMS\SysTrace Agent\SysTrace Agent.lnk" "$INSTDIR\SysTrace_Agent.exe"
    CreateShortcut "$SMPROGRAMS\SysTrace Agent\Deinstallieren.lnk" "$INSTDIR\Uninstall.exe"

    CreateShortCut "$DESKTOP\SysTrace Agent.lnk" "$INSTDIR\SysTrace_Agent.exe"

    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "DisplayName" "SysTrace Agent"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "UninstallString" "$INSTDIR\Uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "InstallLocation" "$INSTDIR"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "DisplayVersion" "1.0.0"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent" "Publisher" "Elias"

SectionEnd

; ============================================
Section "Uninstall"

    ExecWait 'powershell -Command "Get-AppxPackage -Name GpsHelper | Remove-AppxPackage"'

    Delete "$INSTDIR\SysTrace_Agent.exe"
    Delete "$INSTDIR\GpsHelper.cer"
    Delete "$INSTDIR\Uninstall.exe"
    RMDir "$INSTDIR"

    Delete "$DESKTOP\SysTrace Agent.lnk"
    Delete "$SMPROGRAMS\SysTrace Agent\SysTrace Agent.lnk"
    Delete "$SMPROGRAMS\SysTrace Agent\Deinstallieren.lnk"
    RMDir "$SMPROGRAMS\SysTrace Agent"

    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\SysTrace_Agent"

SectionEnd