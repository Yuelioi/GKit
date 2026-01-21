package response

func (r *Response) Gin(c interface {
	JSON(int, interface{})
}) {
	c.JSON(r.Status(), r)
}
