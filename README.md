### Предисловие

Репозиторий создан для изучения работы с числовыми типами в приложениях написанных на go + postgres.
Желательно понимать как работают числа с плавающей точкой под капотом и стандарт IEEE 754  

Полезные ссылки:

[Что нужно знать про арифметику с плавающей запятой](https://habr.com/ru/post/112953/)
[wikipedia IEEE 754](https://ru.wikipedia.org/wiki/IEEE_754-2008)

### Типы данных POSTGRES

[дока](https://postgrespro.ru/docs/postgresql/10/datatype-numeric)

 Числовые типы
| Имя              | Размер     | Описание                                  | Диапазон                                              |
|------------------|------------|-------------------------------------------|-------------------------------------------------------|
| smallint         | 2 байта    | целое в небольшом диапазоне               | -32768 .. +32767                                      |
| integer          | 4 байта    | типичный выбор для целых чисел            | -2147483648 .. +2147483647                            |
| bigint           | 8 байт     | целое в большом диапазоне                 | -9223372036854775808 .. 9223372036854775807           |
| decimal          | переменный | вещественное число с указанной точностью  | до 131072 цифр до десятичной точки и до 16383 — после |
| numeric          | переменный | вещественное число с указанной точностью  | до 131072 цифр до десятичной точки и до 16383 — после |
| real             | 4 байта    | вещественное число с переменной точностью | точность в пределах 6 десятичных цифр                 |
| double precision | 8 байт     | вещественное число с переменной точностью | точность в пределах 15 десятичных цифр                |
| smallserial      | 2 байта    | небольшое целое с автоувеличением         | 1 .. 32767                                            |
| serial           | 4 байта    | целое с автоувеличением                   | 1 .. 2147483647                                       |
| bigserial        | 8 байт     | большое целое с автоувеличением           | 1 .. 9223372036854775807                              |

#### Целочисленные типы.
`integer`, `smallint` и `bigint`

#### Числа с произвольной точностью
`decimal`, `numeric`

`numeric` позволяет сохранить специальное значение NaN.

НО! В большинстве реализаций NaN считается не равным любому другому значению (в том числе и самому NaN. Чтобы значения numeric можно было сортировать и использовать в древовидных индексах, PostgreSQL считает, что значения NaN равны друг другу и при этом больше любых числовых значений (не NaN).

Типы `decimal` и `numeric` равнозначны. Оба эти типа описаны в стандарте SQL.

При округлении значений тип numeric выдаёт число, большее по модулю, тогда как (на большинстве платформ) типы `real` и `double precision` выдают ближайшее чётное число.

```
SELECT x,
  round(x::numeric) AS num_round,
  round(x::double precision) AS dbl_round
FROM generate_series(-3.5, 3.5, 1) as x;
  x   | num_round | dbl_round
------+-----------+-----------
 -3.5 |        -4 |        -4
 -2.5 |        -3 |        -2
 -1.5 |        -2 |        -2
 -0.5 |        -1 |        -0
  0.5 |         1 |         0
  1.5 |         2 |         2
  2.5 |         3 |         2
  3.5 |         4 |         4
(8 rows)
```

####  Типы с плавающей точкой

Типы данных `real` и `double precision` хранят приближённые числовые значения с переменной точностью. На практике эти типы обычно реализуют стандарт IEEE 754 для двоичной арифметики с плавающей точкой (с одинарной и двойной точностью соответственно), в той мере, в какой его поддерживают процессор, операционная система и компилятор.

В дополнение к обычным числовым значениям типы с плавающей точкой могут содержать следующие специальные значения:
Infinity  
-Infinity  
NaN  

PostgreSQL также поддерживает форматы `float` и `float(p)`, оговорённые в стандарте SQL, для указания неточных числовых типов. Здесь `p` определяет минимально допустимую точность в двоичных цифрах. PostgreSQL воспринимает запись от `float(1)` до `float(24)` как выбор типа `real`, а запись от `float(25)` до `float(53)` как выбор типа `double precision`. Значения `p` вне допустимого диапазона вызывают ошибку. Если `float` указывается без точности, подразумевается тип `double precision`.

#### Последовательные типы

Типы данных `smallserial`, `serial` и `bigserial` не являются настоящими типами, а представляют собой просто удобное средство для создания столбцов с уникальными идентификаторами (подобное свойству AUTO_INCREMENT в некоторых СУБД). В текущей реализации запись: 

```
CREATE TABLE имя_таблицы (
    имя_столбца SERIAL
);
```

равнозначна:

```
CREATE SEQUENCE имя_таблицы_имя_столбца_seq AS integer;
CREATE TABLE имя_таблицы (
    имя_столбца integer NOT NULL DEFAULT nextval('имя_таблицы_имя_столбца_seq')
);
ALTER SEQUENCE имя_таблицы_имя_столбца_seq OWNED BY имя_таблицы.имя_столбца;
```

Чтобы в столбец нельзя было вставить `NULL`, в его определение добавляется ограничение `NOT NULL`. (Во многих случаях также имеет смысл добавить для этого столбца ограничения `UNIQUE` или `PRIMARY KEY` для защиты от ошибочного добавления дублирующихся значений, но автоматически это не происходит)

Так как типы `smallserial`, `serial` и `bigserial` реализованы через последовательности, в числовом ряду значений столбца могут образовываться пропуски (или "дыры"), даже если никакие строки не удалялись. Значение, выделенное из последовательности, считается "задействованным", даже если строку с этим значением не удалось вставить в таблицу. Это может произойти, например, при откате транзакции, добавляющей данные.

### Golang

В Go по умолчанию нет типа `decimal` для расчетов с числами с плавающей точкой, поэтому необходимо использовать библиотеку `math/big` и тип `big.Float` например.

Однако же, есть еще библиотека `github.com/shopspring/decimal` которая реализует тип Decimal. Потестим ее.

Трудности - кейсы `c3` и `c11` для `real` отваливаются. Починить это можно, если в запросе переводить в `numeric` - `real_param::numeric` . 

Причина видимо в том что в драйвере `github.com/lib/pq`  типы `real` и `double precision` возвращаются как набор байт в Scan() а не `float64`, как заявлено в доке. Это очень меня расстроило.




