#

```
npm install -g plist2
npm install -g @vscode/vsce
plist2 syntaxes/golden.tmLanguage.yaml syntaxes/golden.tmLanguage.json
vsce package --allow-missing-repository --skip-license --allow-unused-files-pattern
code --install-extension golden-syntax-highlighter-0.0.0.vsix 
```


