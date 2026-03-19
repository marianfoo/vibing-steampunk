CLASS zcl_token DEFINITION PUBLIC FINAL CREATE PUBLIC.
  PUBLIC SECTION.
    METHODS constructor IMPORTING iv_pos TYPE REF TO zcl_position iv_str TYPE string.
    METHODS get_str RETURNING VALUE(rv_result) TYPE string.
    METHODS get_pos RETURNING VALUE(rv_result) TYPE REF TO zcl_position.
  PRIVATE SECTION.
    DATA mo_pos TYPE REF TO zcl_position.
    DATA mv_str TYPE string.
ENDCLASS.

CLASS zcl_token IMPLEMENTATION.
  METHOD constructor.
    me->mv_pos = lv_pos.
    me->mv_str = lv_str.
  ENDMETHOD.
  METHOD get_str.
    rv_result = me->mv_str.
    RETURN.
  ENDMETHOD.
  METHOD get_pos.
    rv_result = me->mv_pos.
    RETURN.
  ENDMETHOD.
ENDCLASS.
