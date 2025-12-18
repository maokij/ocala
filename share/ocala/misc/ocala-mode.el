;;; ocala-mode.el -- Major mode for editing Ocala code.

;;; Commentary:

;; init.el:
;; (use-package ocala-mode
;;   :mode "\\.oc\\'"
;;   :config (add-hook 'ocala-mode-hook #'eglot-ensure))

;;; Code:
(require 'smie)
(require 'rx)

(defvar ocala-indent-basic 4)

(defvar ocala-mode-map
  (let ((map (make-sparse-keymap)))
    map)
  "Keymap for `ocala-mode'.")

(defconst ocala-mode-syntax-table
  (let ((table (make-syntax-table)))
    (modify-syntax-entry ?\' "\"" table)
    (modify-syntax-entry ?$ "_" table)
    (modify-syntax-entry ?% "_" table)
    (modify-syntax-entry ?& "_" table)
    (modify-syntax-entry ?+ "_" table)
    (modify-syntax-entry ?- "_" table)
    (modify-syntax-entry ?/ "_ 124b" table)
    (modify-syntax-entry ?* "_ 23" table)
    (modify-syntax-entry ?< "_" table)
    (modify-syntax-entry ?= "_" table)
    (modify-syntax-entry ?> "_" table)
    (modify-syntax-entry ?| "_" table)
    (modify-syntax-entry ?? "_" table)
    (modify-syntax-entry ?! "_" table)
    (modify-syntax-entry ?{ "(}" table)
    (modify-syntax-entry ?} "){" table)
    (modify-syntax-entry ?\[ "(]" table)
    (modify-syntax-entry ?\] ")[" table)
    (modify-syntax-entry ?\( "()" table)
    (modify-syntax-entry ?\) ")(" table)
    (modify-syntax-entry ?\n "> b" table)
    table)
  "Syntax table for `ocala-mode`.")

(defconst ocala-registers
  '("A" "B" "C" "D" "E" "H" "L" "AF" "BC" "DE" "HL" "AF-"
    "IX" "IY" "IXH" "IXL" "IYH" "IYL"
    "X" "Y" "S" "P"
    "SP" "PC" "__PC__" "__PROC__"))

(defconst ocala-conditions
  '("NZ?" "NE?" "not-zero?" "!=?"
    "Z?" "EQ?" "zero?" "==?"
    "NC?" "CC?" "not-carry?" "borrow?" ">=?"
    "C?" "CS?" "carry?" "not-borrow?" "<?"
    "PO?" "VC?" "odd?" "not-over?"
    "PE?" "VS?" "even?" "over?"
    "P" "PL?" "plus?"
    "M" "MI?" "minus?"))

(defconst ocala-reserved-words
  '("*patch*"
    "<reserved>"
    "alias"
    "align"
    "apply"
    "arch"
    "assert"
    "break"
    "break-if"
    "case"
    "compile-error"
    "const"
    "continue"
    "continue-if"
    "data"
    "debug-inspect"
    "defined?"
    "do"
    "else"
    "expand-loop"
    "exprtypeof"
    "fallthrough"
    "flat!"
    "goto"
    "goto-if"
    "if"
    "import"
    "incbin"
    "include"
    "link"
    "load-file"
    "loop"
    "macro"
    "make-counter"
    "module"
    "nameof"
    "nametypeof"
    "never-return"
    "once"
    "optimize"
    "pragma"
    "proc"
    "quote"
    "recur"
    "redo"
    "redo-if"
    "return"
    "return-if"
    "section"
    "sizeof"
    "struct"
    "tco"
    "warn"
    "when"))

(defconst ocala-binary-operators
  '("<-"
    "->"
    "<->"
    "+"
    "+$"
    "&"
    "|"
    "^"
    "<<"
    ">>"
    ">>>"
    "*"
    "/"
    "%"
    "<"
    "<="
    ">"
    ">="
    "=="
    "!="
    "-set"
    "-reset"
    "-bit?"
    "-in"
    "-out"
    "-$"
    "-?"
    "-"
    "<*"
    "<*$"
    ">*"
    ">*$"
    "<<"
    ">>"
    "<<<"
    ">>>"))

(defconst ocala-operators-tail-re "\\([^({[]\\|$\\)")

(defconst ocala-binary-operators-re
  (rx (or "." (regexp (regexp-opt ocala-binary-operators 'symbols)))
      (regexp ocala-operators-tail-re)))

(defconst ocala-unary-operators
  '("++"
    "--"
    "-not"
    "-neg"
    "-push"
    "-pop"
    "-zero?"))

(defconst ocala-unary-operators-re
  (rx (regexp (regexp-opt ocala-unary-operators 'symbols))
      (regexp ocala-operators-tail-re)))

(defconst ocala-operators-re
  (rx (or "." (regexp (regexp-opt (append ocala-binary-operators ocala-unary-operators) 'symbols)))
      (regexp ocala-operators-tail-re)))

(defconst ocala-font-lock-keywords
  (list
   (cons
    (regexp-opt (append ocala-registers ocala-reserved-words) 'symbols)
    'font-lock-keyword-face)
   (cons
    (regexp-opt (append ocala-conditions ocala-binary-operators ocala-unary-operators) 'symbols)
    'font-lock-function-name-face)))

(defconst ocala-smie-grammar
  (smie-prec2->grammar
   (smie-bnf->prec2
    '((id)
      (insts (insts ";" insts) (inst))
      (inst ("{" insts "}")
            ("(" insts ")")
            ("[" insts "]")
            (exp))
      (exp (exp "," exp)
           (exp "BOP" exp)
           (exp "UOP")))
    '((assoc ";"))
    '((assoc ",") (assoc "BOP") (assoc "UOP")))))

(defun ocala-smie--implicit-semi-p ()
  (save-excursion
    (skip-chars-backward " \t")
    (not (or (memq (char-before) '(?\{ ?\[ ?\( ?\,))
             (looking-back ocala-binary-operators-re (- (point) 8) t)
             (forward-comment (point-max))
             (looking-at ocala-operators-re)))))

(defun ocala-smie--multi-line-seq-p ()
  (save-excursion
    (skip-chars-backward " \t")
    (let* ((p (nth 1 (syntax-ppss)))
           (q (save-excursion (re-search-backward "[^,][ \t]*(//.*)?$" p t))))
      (and q (setq p (save-excursion
                       (goto-char q)
                       (end-of-line)
                       (point))))
      (re-search-backward ", *//-" p t))))

(defun ocala-smie--forward-token ()
  (let ((pos (point)))
    (skip-chars-forward " \t")
    (cond
     ((and (looking-at (rx-to-string `(| "\n" "//" "/*")))
           (ocala-smie--implicit-semi-p))
      (if (eolp) (forward-char 1) (forward-comment 1))
      ";")
     (t
      (forward-comment (point-max))
      (cond
       ((looking-at "[]})({[;]")
        (forward-char 1)
        (match-string 0))
       ((looking-at ",")
        (forward-char 1)
        (if (or (looking-at " *//-")
                (ocala-smie--multi-line-seq-p))
            "BOP" ","))
       ((looking-at "\\s\"")
        (forward-sexp 1)
        "STR")
       ((looking-at ocala-binary-operators-re)
        (goto-char (match-end 0))
        "BOP")
       ((looking-at ocala-unary-operators-re)
        (goto-char (match-end 0))
        "UOP")
       (t (smie-default-forward-token)))))))

(defun ocala-smie--backward-token ()
  (let ((pos (point)))
    (forward-comment (- (point)))
    (cond
     ((and (> pos (line-end-position))
           (ocala-smie--implicit-semi-p))
      ";")
     ((looking-back "[]})({[;]" (- (point) 1) t)
      (forward-char -1)
      (match-string 0))
     ((looking-back "," (- (point) 1))
      (forward-char -1)
      (if (or (looking-at ", *//-")
              (ocala-smie--multi-line-seq-p))
          "BOP" ","))
     ((looking-back "\\s\"" (- (point) 2) t)
      (backward-sexp 1)
      "STR")
     ((looking-back ocala-binary-operators-re (- (point) 8) t)
      (goto-char (match-beginning 0))
      "BOP")
     ((looking-back ocala-unary-operators-re (- (point) 8) t)
      (goto-char (match-beginning 0))
      "UOP")
     (t (smie-default-backward-token)))))

(defun ocala-smie-rules (kind token)
  (pcase (cons kind token)
    (`(:elem . basic) ocala-indent-basic)
    (`(,_ . ",") (smie-rule-separator kind))
    (`(:after . "BOP") (and (smie-rule-sibling-p) (- ocala-indent-basic)))
    (`(:before . ,(or `"BOP" `"UOP"))
     (cons 'column
           (save-excursion
             (forward-comment (- (point)))
             (let* ((a (progn (beginning-of-line) (point)))
                    (b (progn (skip-chars-forward " \t") (point)))
                    (c (- b a)))
               (if (looking-at ocala-operators-re)
                   c
                 (+ c ocala-indent-basic))))))
    (`(:before . ,(or `"[" `"(" `"{"))
     (unless (looking-back "^[ \t]+")
       (smie-rule-parent)))))

(defun ocala-imenu-create-index ()
  (let ((index-alist '())
        (end (point-max))
        name pos decl)
    (goto-char (point-min))
    (while (re-search-forward "^\\s *\\(macro\\|proc\\)\\s +\\([^(\n ]+\\)" end t)
      (setq decl (match-string-no-properties 1))
      (setq name (match-string-no-properties 2))
      (setq pos (match-beginning 0))
      (push (cons name pos) index-alist)
      (re-search-forward "{")
      (smie-forward-sexp))
    (nreverse index-alist)))

;;;###autoload
(define-derived-mode ocala-mode prog-mode "Ocala"
  "Major mode for editing Ocala code."
  (set-syntax-table ocala-mode-syntax-table)
  (setq-local font-lock-defaults '(ocala-font-lock-keywords nil nil))
  (setq-local imenu-create-index-function #'ocala-imenu-create-index)

  (setq-local comment-start "// ")
  (setq-local comment-end "")
  (setq-local comment-start-skip "\\(//+\\|/\\*+\\)\\s *")
  (setq-local comment-use-syntax t)
  (setq-local comment-multi-line t)

  (smie-setup ocala-smie-grammar #'ocala-smie-rules
              :forward-token #'ocala-smie--forward-token
              :backward-token #'ocala-smie--backward-token)
  (use-local-map ocala-mode-map))

;;;###autoload
(add-to-list 'auto-mode-alist '("\\.oc\\'" . ocala-mode))

(with-eval-after-load 'eglot
  (add-to-list 'eglot-server-programs '((ocala-mode) "ocala-language-server")))

(provide 'ocala-mode)

;;; ocala-mode.el ends here
