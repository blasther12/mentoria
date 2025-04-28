- ✅ Inicia duas operações (DB e HTTP) em paralelo.
- ✅ Se alguma delas demorar mais de 500ms, cancela tudo.
- ✅ Se o usuário apertar CTRL+C, cancela com segurança.
- ✅ Espera as duas operações terminarem para finalizar.
- ✅ Protege variáveis compartilhadas contra condição de corrida.


```
[Início main()]
      |
      v
[Cria Context com Timeout de 500ms]
      |
      v
[Escuta sinais SIGINT/SIGTERM]
      |
      v
[Inicia goroutine runOperations(ctx)]
      |
      v
[runOperations]
    ├──> [Goroutine dbConnection(ctx)]
    └──> [Goroutine httpConnection(ctx)]
      |
      v
[WAIT: Espera dbConnection + httpConnection terminarem]
      |
      v
---------------------------------------------
SELECIONA:
   (caso 1) <- timeout 500ms?
       |
       v
  [Context canceled -> operações recebem <-ctx.Done()]
       |
       v
  [dbConnection ou httpConnection abortam e retornam erro]

OU

   (caso 2) <- recebimento de sinal (ex: CTRL+C)?
       |
       v
  [Cancela Context manualmente -> operações cancelam]

OU

   (caso 3) <- operações terminam normalmente antes dos 500ms?
       |
       v
  [Termina com sucesso]
---------------------------------------------
      |
      v
[Fecha canal done]
      |
      v
[main() detecta fim e imprime: Operações finalizadas OU Erro ocorrido]
      |
      v
[Fim]

```