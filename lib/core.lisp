(ns 'types)

(def Int    (type 0))
(def Float  (type 0.0))
(def Vector (type []))
(def List   (type ()))
(def Set    (type #{}))
(def Bool   (type true))
(def String (type "specimen"))
(def Keyword(type :specimen))
(def Symbol (type 'specimen))

(ns 'core)

(def fn
    (macro* fn
        ([& decl] (cons 'fn* decl))))

(def defn (macro* defn [name & fdecl]
    (let* [func (cons 'fn* (cons name fdecl))]
    `(def ~name ~func))))

(def defmacro (macro* defmacro [name & mdecl]
    (let* [macro (cons 'macro* (cons name mdecl))]
    `(def ~name ~macro))))

(defn nil? [arg] (= nil arg))

(defn empty? [coll]
    (if (nil? coll)
        true
        (nil? (first coll))))

; Type check functions -------------------------------
(defn is-type? [typ arg] (= typ (type arg)))
(defn seq? [arg] (impl? arg types/Seq))
(defn set? [arg] (is-type? #{} arg))
(defn list? [arg] (is-type? types/List arg))
(defn vector? [arg] (is-type? types/Vector arg))
(defn int? [arg] (is-type? types/Int arg))
(defn float? [arg] (is-type? types/Float arg))
(defn boolean? [arg] (is-type? types/Bool arg))
(defn string? [arg] (is-type? types/String arg))
(defn keyword? [arg] (is-type? types/Keyword arg))
(defn symbol? [arg] (is-type? types/Symbol arg))

; Type initialization functions ---------------------
(defn set [coll] (apply-seq (type #{}) coll))
(defn list [& coll] (apply-seq (type ()) coll))
(defn vector [& coll] (apply-seq (type []) coll))
(defn int [arg] (to-type arg (type 0)))
(defn float [arg] (to-type arg (type 0.0)))
(defn boolean [arg] (true? arg))

(defn number? [num]
    (if (float? num)
        true
        (int? num)))


; boolean operations --------------------------------
(defn true? [arg]
    (if (nil? arg)
        false
        (if (boolean? arg)
            arg
            true)))

(defn not [arg] (= false (true? arg)))


; sequence operations -------------------------------
(defn cons [val coll]
    (if (nil? coll)
        (cons val ())
        (if (seq? coll)
            ((. Cons coll) val)
            (throw "cons cannot be done for " (type coll)))))

(defn last [coll]
    (let* [v   (first coll)
           rem (next coll)]
        (if (nil? rem)
            v
            (last (next coll)))))
