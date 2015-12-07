;; add this to go-mode-hook
;; (add-hook 'go-mode-hook 'goautoenv-path-setup)

(defun goautoenv-path-setup
    (defun append-path (a &rest arg-list)
      (dolist (b arg-list)
        (setq a (concat (file-name-as-directory a) b)))
      a)
  (setq path (file-name-directory buffer-file-name))
  (while path
    (setq tmp_path (append-path path ".goenv" "bin" "activate"))
    (if (file-exists-p tmp_path)
        (progn
          (make-local-variable 'goautoenv)
          (setq goautoenv
                (with-temp-buffer
                  (insert-file-contents tmp_path)
                  (substring (buffer-string)
                             (+ 7 (string-match "^GOPATH=.+$" (buffer-string)))
                             (match-end 0))))
          (add-function :before (symbol-function 'company-go--invoke-autocomplete) #'goautoenv-invoke-autocomplete)
          (setq path nil))
      (setq path (file-name-directory (directory-file-name path))))))
