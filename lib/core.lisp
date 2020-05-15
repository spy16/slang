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
(def Fn     (type (fn* [])))

(ns 'core)

(def fn (macro* fn [& decl]
    (decl.Cons 'fn*)))

(def defn (macro* defn [name & fdecl]
    (let* [with-name (fdecl.Cons name)
           func      (fdecl.Cons 'fn*)]
        `(def ~name ~func))))

(def defmacro (macro* defmacro [name & mdecl]
    (let* [with-name (mdecl.Cons name)
           macro     (mdecl.Cons 'macro*)]
        `(def ~name ~macro))))

(defn nil? [arg] (= nil arg))

; sequence operations -------------------------------
(defn seq? [arg] (impl? arg types/Seq))

(defn first [coll]
    (if (nil? coll)
        nil
        (if (not (seq? coll))
            (throw "argument must be a collection, not " (type coll))
            (coll.First))))

(defn second [coll]
    (first (next coll)))

(defn next [coll]
    (if (not (seq? coll))
        (throw "argument must be a collection, not " (type coll)))
    (coll.Next))

(defn cons [v coll]
    (if (not (seq? coll))
        (throw "argument must be a collection, not " (type coll)))
    (coll.Cons v))

(defn conj [coll & vals]
    (if (not (seq? coll))
        (throw "argument must be a collection, not " (type coll)))
    (apply-seq coll.Conj vals))

(defn empty? [coll]
    (if (nil? coll)
        true
        (nil? (first coll))))

(defn cons [val coll]
    (if (nil? coll)
        (cons val ())
        (if (seq? coll)
            (coll.Cons val)
            (throw "cons cannot be done for " (type coll)))))

(defn last [coll]
    (let* [v   (first coll)
           rem (next coll)]
        (if (nil? rem)
            v
            (last (next coll)))))

(defn number? [num]
    (if (float? num)
        true
        (int? num)))

(defn even? [num]
    (= (mod num 2) 0.0))

(defn odd? [num]
    (= (mod num 2) 1.0))

(defn inc [num]
    (if (not (number? num))
        (throw "argument must be a number"))
    (if (int? num)
        (int (+ 1 num))
        (+ 1 num)))

(defn count
    ([coll] (count coll 0))
    ([coll counter]
        (if (empty? coll)
            counter
            (count (next coll) (inc counter)))))

(defn concat [coll1 coll2]
    (if (not (seq? coll1))
        (throw "argument coll1 must be a seq, not " (type coll1)))
    (if (not (seq? coll2))
        (throw "argument coll2 must be a seq, not " (type coll2)))
    (apply-seq coll1.Conj coll2))

(defn map
    ([f coll]
        (if (empty? coll)
            nil)
        (map f coll []))
    ([f coll acc]
        (if (empty? coll)
            acc
            (map f (next coll) (conj acc (f (first coll)))))))

(defn filter
    ([f coll]
        (if (empty? coll)
            nil)
        (filter f coll []))
    ([f coll acc]
        (if (empty? coll)
            acc
            (if (f (first coll))
                (filter f (next coll) (conj acc (first coll)))
                (filter f (next coll) acc)))))

; important macros -----------------------------------
(defmacro apply-seq [callable args]
    `(eval (cons ~callable ~args)))

(defmacro when [expr & body]
    (let* [body (cons 'do body)]
    `(if ~expr ~body)))

(defmacro when-not [expr & body]
    (let* [body (cons 'do body)]
    `(if (not ~expr) ~body)))

(defmacro assert
    ([expr] (let* [message "assertion failed"]
                `(when-not ~expr (throw ~message))))
    ([expr message] `(when-not ~expr (throw ~message))))

; Type check functions -------------------------------
(defn is-type? [typ arg] (= typ (type arg)))
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

; boolean operations --------------------------------
(defn true? [arg]
    (if (nil? arg)
        false
        (if (boolean? arg)
            arg
            true)))

(defn not [arg] (= false (true? arg)))
