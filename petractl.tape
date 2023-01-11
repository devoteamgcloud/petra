Output petractl.gif

Set FontSize 22
Set Width 1200
Set Height 600

Type "petractl"
Sleep 100ms
Enter

Sleep 1s

Ctrl+L

Sleep 1s

Type "cat ./tests/mod1/main.tf"

Sleep 1s 

Ctrl+L

Type "cat ./tests/mod1/petra-config.yaml"

Sleep 1s

Ctrl+L

Type "petractl push --bucket orga-automation-petra-backend-30ba ./tests/mod1"

Sleep 2s