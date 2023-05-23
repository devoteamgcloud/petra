Output petractl.gif

Set FontSize 21
Set WindowBar "Colorful"
Set Padding 2
Set Width 1200
Set Height 600
Set PlaybackSpeed 0.8

Type "petractl"
Sleep 100ms
Enter

Sleep 2s

Ctrl+L

Sleep 1s

Type "cat ./tests/mod1/main.tf"
Sleep 100ms
Enter

Sleep 2s 

Ctrl+L

Type "cat ./tests/mod1/petra-config.yaml"
Sleep 100ms
Enter

Sleep 2s

Ctrl+L

Type "petractl push --bucket tf-registry-petra ./tests/mod1"
Sleep 100ms
Enter
Sleep 3s