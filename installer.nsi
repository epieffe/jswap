# This is a NSIS script to generate the Jswap installer for Windows.
# For more info visit https://nsis.sourceforge.io/

!include winmessages.nsh
!include LogicLib.nsh

# Do not request admin privileges
RequestExecutionLevel user

Name Jswap
Outfile "jswap-setup.exe"

# By default Jswap is installed in the same folder where it stores data,
# but the user might change the installation directory during setup.
InstallDir $LocalAppData\Jswap

!define DATAFOLDER $LocalAppData\Jswap
!define JAVAHOME ${DATAFOLDER}\current-jdk

# License page
LicenseData LICENSE
Page license

Function checkJavaHome
    # Warn the user if JAVA_HOME environment variable is already set
    ReadEnvStr $0 "JAVA_HOME"
    StrCmp $0 "" continue
    MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "Warning: the JAVA_HOME environment variable is already set on your system, installing Jswap will override it." IDOK continue IDCANCEL abort
    abort:
        Quit
    continue:
FunctionEnd

# Installation directory selection page
Page directory checkJavaHome

# Installation page where sections are executed
Page instfiles

Section "Install"

    # Copy jswap.exe to install folder
    SetOutPath $INSTDIR\bin
    File jswap.exe

    # Add jswap to PATH, only if not present
    EnVar::AddValue "Path" "$INSTDIR\bin"
    Pop $0
    DetailPrint "Added Jswap to PATH. result=$0"

    # Add JDK symlink to PATH, only if not present
    EnVar::AddValue "Path" "${JAVAHOME}\bin"
    Pop $0
    DetailPrint "Added JDK symlink to PATH. result=$0"

    # Set JAVA_HOME environment variable
    WriteRegExpandStr HKCU "Environment" "JAVA_HOME" "${JAVAHOME}"
    DetailPrint "Set JAVA_HOME environment variable"

    # Make sure Windows knows about the environment change
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

    # Create uninstaller
    WriteUninstaller $INSTDIR\uninstall.exe

SectionEnd


# This section defines what the uninstaller does
Section "Uninstall"

    # Delete jswap from PATH, only if present
    EnVar::DeleteValue "Path" "$INSTDIR\bin"
    Pop $0
    DetailPrint "Removed Jswap from PATH. result=$0"

    # Delete JDK symlink from PATH, only if present
    EnVar::DeleteValue "Path" "${JAVAHOME}\bin"
    Pop $0
    DetailPrint "Removed JDK symlink from PATH. result=$0"

    # Delete JAVA_HOME environment variable, only if it points to JDK symlink
    ReadEnvStr $0 "JAVA_HOME"
    ${If} $0 == ${JAVAHOME}
        DeleteRegValue HKCU "Environment" "JAVA_HOME"
        DetailPrint "Removed JAVA_HOME environment variable"
    ${Else}
        DetailPrint "JAVA_HOME environment variable was not set"
    ${EndIf}

    # Make sure Windows knows about the environment change
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

    # Delete installation files
    Delete /REBOOTOK $INSTDIR\bin\jswap.exe
    Delete /REBOOTOK $INSTDIR\uninstall.exe
    RMDir /REBOOTOK $INSTDIR\bin
    # Delete Jswap data folder
    RMDir /r /REBOOTOK ${DATAFOLDER}
    # Delete installation folder, only if it is empty
    RMDir /REBOOTOK $INSTDIR

SectionEnd
