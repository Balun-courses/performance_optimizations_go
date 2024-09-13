## Структура процессора
1) Память команд и данных, модуль [memory.v](./cpu-template/memory.v)
2) Регистровый файл, модуль [register_file.v](./cpu-template/register_file.v)
3) Модуль для одиночного 32-битного регистра PC в [d_flop.v](./cpu-template/d_flop.v)
4) Модуль [util.v](./cpu-template/util.v) содержит вспомогательные модули, которые могут быть полезны
   при реализации процессора.


| Команда | opcode | rs    | rt    | imm              |
|---------|--------|-------|-------|------------------|
| lw      | 100011 | xxxxx | xxxxx | xxxxxxxxxxxxxxxx |
   | sw      | 101011 | xxxxx | xxxxx | xxxxxxxxxxxxxxxx |
| beq     | 000100 | xxxxx | xxxxx | xxxxxxxxxxxxxxxx |

| Команда | opcode | rs    | rt    | rd    | shamt | funct  |
|---------|--------|-------|-------|-------|-------|--------|
| add     | 000000 | xxxxx | xxxxx | xxxxx | 00000 | 100000 |
| sub     | 000000 | xxxxx | xxxxx | xxxxx | 00000 | 100010 |
| and     | 000000 | xxxxx | xxxxx | xxxxx | 00000 | 100100 |
| or      | 000000 | xxxxx | xxxxx | xxxxx | 00000 | 100101 |
| slt     | 000000 | xxxxx | xxxxx | xxxxx | 00000 | 101010 |

| Команда | opcode | rs    | rt    | imm              |
|---------|--------|-------|-------|------------------|
| addi*   | 001000 | xxxxx | xxxxx | xxxxxxxxxxxxxxxx |
| andi*   | 001100 | xxxxx | xxxxx | xxxxxxxxxxxxxxxx |
| bne     | 000101 | xxxxx | xxxxx | xxxxxxxxxxxxxxxx |

| Команда | opcode | addr                       |
|---------|--------|----------------------------|
| j       | 000010 | xxxxxxxxxxxxxxxxxxxxxxxxxx |
| jal     | 000011 | xxxxxxxxxxxxxxxxxxxxxxxxxx |
