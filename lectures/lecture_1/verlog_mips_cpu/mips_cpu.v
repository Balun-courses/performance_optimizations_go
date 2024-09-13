`include "util.v"

/*
               NAMING RULES:
- Modules naiming:
    Sample: my_name_name --> my_not_xor_gate(may be "gate"), if it's a helper method
    Sample: name_name --> mips_cpu, mostly for basic methods

- Input and output variables in modules:
    Sample input: number + Input --> firstInput, secondInput, thirdInput...
    Sample output: output + Name --> outputResult, outputValue

- Supply (1/0):
    Sample ground(0): gnd
    Sample voltage drain drain(1): vdd

- Other variables in modules:
    Sample: name_name_name --> a0_and_b0

- Exception to the rule:
    Tested modules and edge cases

N.B. Each module must have spaces between logical blocks of code. Can have comments
*/

module mips_cpu(clk, pc, pc_new, instruction_memory_a, instruction_memory_rd, data_memory_a, data_memory_rd, data_memory_we, data_memory_wd,
    register_a1, register_a2, register_a3, register_we3, register_wd3, register_rd1, register_rd2);
    input clk;
    inout[31:0] pc;
    output[31:0] pc_new;
    output data_memory_we;
    output[31:0] instruction_memory_a, data_memory_a, data_memory_wd;
    inout[31:0] instruction_memory_rd, data_memory_rd;
    output register_we3;
    output[4:0] register_a1, register_a2, register_a3;
    output[31:0] register_wd3;
    inout[31:0] register_rd1, register_rd2;

    // contol unit flag initialization
    wire[5:0] aluContol;

    wire memToReg, memWrite, branch, aluSource, registerDestination,
         registerWrite, branchIfNotEqual, jump, jumpAndLink, jumpRegister;

    control_unit my_mips_control_unit(
        instruction_memory_rd[31:26],
        instruction_memory_rd[5:0],
        memToReg,
        memWrite,
        branch,
        aluContol,
        aluSource,
        registerDestination,
        registerWrite,
        branchIfNotEqual,
        jump,
        jumpAndLink,
        jumpRegister
    );

    // beginning of the circuit, connection of the instruction_memory input
    assign instruction_memory_a = pc;

    // initialization of sign extend const
    wire[31:0] const_sign_extend_value;

    sign_extend util_sign_extend(instruction_memory_rd[15:0], const_sign_extend_value);

    // connection of the register_file input (+jal)
    wire[4:0] register_to_write;

    mux2_5 first_mux(register_to_write, 5'b11111, jumpAndLink, register_a3);
    mux2_5 second_mux(instruction_memory_rd[20:16], instruction_memory_rd[15:11], registerDestination, register_to_write);

    assign register_a1 = instruction_memory_rd[25:21];
    assign register_a2 = instruction_memory_rd[20:16];
    assign register_we3 = registerWrite;


    // initialization of sources from registers (sources to ALU)
    wire[31:0] source_of_first_register_a1, source_of_second_register_a2;
    assign source_of_first_register_a1 = register_rd1;
    mux2_32 mux_for_second_source(register_rd2, const_sign_extend_value, aluSource, source_of_second_register_a2);

    // connection to ALU
    wire zero_flag;
    wire[31:0] alu_result;

    my_alu_for_thirty_two_bit_numbers main_alu(source_of_first_register_a1, source_of_second_register_a2,
        aluContol, alu_result, zero_flag);

    // connection of the instruction_memory input
    assign data_memory_a = alu_result;
    assign data_memory_wd = register_rd2;
    assign data_memory_we = memWrite;

    // results processing (ALU and data memory)
    wire[31:0] result;
    mux2_32 mux_for_alu_and_read_data(alu_result, data_memory_rd, memToReg, result);
    mux2_32 mux_for_write_data_register(result, PCPlus4, jumpAndLink, register_wd3);

    // circuit iteration initialization
    wire[31:0] PCPlus4, PCBranch, bif_shift;
    shl_2 bit_shift(const_sign_extend_value, bif_shift);

    adder adder_to_pc(pc, 4, PCPlus4);
    adder add_to_branch(PCPlus4, bif_shift, PCBranch);

    // jump commands -->
    wire[31:0] jump_addresses, first_jump_command_group, second_jump_command_group;
    assign jump_addresses[1:0] = 2'b00;
    assign jump_addresses[31:28] = 4'b0000;
    assign jump_addresses[27:2] = instruction_memory_rd[25:0];

    // simple jump or jump on alu address
    mux2_32 mux_for_addresses(jump_addresses, source_of_first_register_a1, jumpRegister, second_jump_command_group);

    // jump if branch or branch if not equal
    mux2_32 mux_for_branch_and_branch_if_not_equal(PCPlus4, PCBranch,
    (zero_flag && branch) || (!zero_flag && branchIfNotEqual),
        first_jump_command_group);

    // make a jump or go sequentially
    mux2_32 result_jump_mux(first_jump_command_group, second_jump_command_group, jump, pc_new);
endmodule

module control_unit(operand, currentFunction, memToReg, memWrite, branch, aluControl, aluSource,
    registerDestination, registerWrite, branchIfNotEqual, jump, jumpAndLink, jumpRegister);

    inout[5:0] operand, currentFunction;

    output reg[5:0] aluControl;
    output reg memToReg, memWrite, branch, aluSource,
               registerDestination, registerWrite, branchIfNotEqual, jump, jumpAndLink, jumpRegister;

    always @(*) begin
        memToReg = 0;
        memWrite = 0;
        branch = 0;
        registerDestination = 0;
        branchIfNotEqual = 0;

        jump = 0;
        jumpAndLink = 0;
        jumpRegister = 0;

        case (operand)

            // functions and jump register
            6'b000000: begin
                aluSource = 0;

                if (currentFunction == 6'b001000) begin
                    aluControl = 6'b100000;
                    registerWrite = 0;

                    jumpRegister = 1;
                    jump = 1;
                end else begin
                    aluControl = currentFunction;
                    registerWrite = 1;

                    registerDestination = 1;
                end
            end

            // load word
            6'b100011: begin
                aluSource = 1;
                aluControl = 6'b100000;
                registerWrite = 1;

                memToReg = 1;
            end

            // set word
            6'b101011: begin
                aluSource = 1;
                aluControl = 6'b100000;
                registerWrite = 0;

                memWrite = 1;
            end

            // branch
            6'b000100: begin
                aluSource = 0;
                aluControl = 6'b100010;
                registerWrite = 0;

                branch = 1;
            end

            // branch if not equal
            6'b000101: begin // bne
                aluSource = 0;
                aluControl = 6'b100010;
                registerWrite = 0;

                branchIfNotEqual = 1;
            end

            // and immediate
            6'b001100: begin
                aluSource = 1;
                aluControl = 6'b100100;
                registerWrite = 1;
            end

            // add immediate
            6'b001000: begin
                aluSource = 1;
                aluControl = 6'b100000;
                registerWrite = 1;
            end

            // jump
            6'b000010: begin
                aluSource = 0;
                aluControl = 6'b100000;
                registerWrite = 0;

                jump = 1;
            end

            // jump and link
            6'b000011: begin
                aluSource = 0;
                aluControl = 6'b100000;
                registerWrite = 1;

                jump = 1;
                jumpAndLink = 1;
            end

        endcase
    end
endmodule

module my_alu_for_thirty_two_bit_numbers(firstOperand, secondOperand, controlSignal, aluResult, zeroFlag);
    input[5:0] controlSignal; // Control signals for operation selection
    input[31:0] firstOperand, secondOperand; // Operands

    output zeroFlag; // For equal operands
    output[31:0] aluResult; // Result

    reg reg_zero_flag;
    reg[31:0] reg_alu_result;

    assign zeroFlag = reg_zero_flag;
    assign aluResult = reg_alu_result;

    always @(*) begin
        case (controlSignal)

            6'b100100: begin
                reg_zero_flag = 0;
                reg_alu_result = firstOperand & secondOperand;
            end

            6'b100101: begin
                reg_zero_flag = 0;
                reg_alu_result = firstOperand | secondOperand;
            end

            6'b100000: begin
                reg_zero_flag = 0;
                reg_alu_result = firstOperand+secondOperand;
            end

            6'b100010: begin
                reg_alu_result = firstOperand+~secondOperand+1;

                // Operands are equal
                if (reg_alu_result == 0) begin
                    reg_zero_flag = 1;
                end else begin
                    reg_zero_flag = 0;
                end
            end

            6'b101010: begin
                reg_zero_flag = 0;

                if (firstOperand[31] == secondOperand[31] && firstOperand < secondOperand) begin
                    reg_alu_result = 1;
                end else if (firstOperand[31] == secondOperand[31] && firstOperand >= secondOperand) begin
                    reg_alu_result = 0;
                end else if (firstOperand[31] == 1 && secondOperand[31] == 0) begin
                    reg_alu_result = 1;
                end else if (firstOperand[31] == 0 && secondOperand[31] == 1) begin
                    reg_alu_result = 0;
                end
            end

        endcase
    end
endmodule
