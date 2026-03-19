CLASS zcl_position DEFINITION PUBLIC FINAL CREATE PUBLIC.
  PUBLIC SECTION.
    METHODS constructor IMPORTING iv_row TYPE i iv_col TYPE i.
    METHODS get_row RETURNING VALUE(rv_result) TYPE i.
    METHODS get_col RETURNING VALUE(rv_result) TYPE i.
  PRIVATE SECTION.
    DATA mv_row TYPE i.
    DATA mv_col TYPE i.
ENDCLASS.

CLASS zcl_position IMPLEMENTATION.
  METHOD constructor.
    me->mv_row = lv_row.
    me->mv_col = lv_col.
  ENDMETHOD.
  METHOD get_row.
    rv_result = me->mv_row.
    RETURN.
  ENDMETHOD.
  METHOD get_col.
    rv_result = me->mv_col.
    RETURN.
  ENDMETHOD.
ENDCLASS.
